// Copyright 2016 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package restorable

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2/internal/graphics"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicscommand"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver"
)

type Pixels struct {
	pixelsRecords *pixelsRecords
}

// Apply applies the Pixels state to the given image especially for restoring.
func (p *Pixels) Apply(img *graphicscommand.Image) {
	// Pixels doesn't clear the image. This is a caller's responsibility.

	if p.pixelsRecords == nil {
		return
	}
	p.pixelsRecords.apply(img)
}

func (p *Pixels) AddOrReplace(pix *graphics.ManagedBytes, region image.Rectangle) {
	if p.pixelsRecords == nil {
		p.pixelsRecords = &pixelsRecords{}
	}
	p.pixelsRecords.addOrReplace(pix, region)
}

func (p *Pixels) Clear(region image.Rectangle) {
	// Note that we don't care whether the region is actually removed or not here. There is an actual case that
	// the region is allocated but nothing is rendered. See TestDisposeImmediately at shareable package.
	if p.pixelsRecords == nil {
		return
	}
	p.pixelsRecords.clear(region)
}

func (p *Pixels) ReadPixels(pixels []byte, region image.Rectangle, imageWidth, imageHeight int) {
	if p.pixelsRecords == nil {
		for i := range pixels {
			pixels[i] = 0
		}
		return
	}
	p.pixelsRecords.readPixels(pixels, region, imageWidth, imageHeight)
}

func (p *Pixels) AppendRegion(regions []image.Rectangle) []image.Rectangle {
	if p.pixelsRecords == nil {
		return regions
	}
	return p.pixelsRecords.appendRegions(regions)
}

func (p *Pixels) Dispose() {
	if p.pixelsRecords == nil {
		return
	}
	p.pixelsRecords.dispose()
}

type ImageType int

const (
	// ImageTypeRegular indicates the image is a regular image.
	ImageTypeRegular ImageType = iota

	// ImageTypeScreen indicates the image is used as an actual screen.
	ImageTypeScreen

	// ImageTypeVolatile indicates the image is cleared whenever a frame starts.
	//
	// Regular non-volatile images need to record drawing history or read its pixels from GPU if necessary so that all
	// the images can be restored automatically from the context lost. However, such recording the drawing history or
	// reading pixels from GPU are expensive operations. Volatile images can skip such operations, but the image content
	// is cleared every frame instead.
	ImageTypeVolatile
)

// Image represents an image that can be restored when GL context is lost.
type Image struct {
	image *graphicscommand.Image

	width  int
	height int

	basePixels Pixels

	// stale indicates whether the image needs to be synced with GPU as soon as possible.
	stale bool

	// staleRegions indicates the regions to restore.
	// staleRegions is valid only when stale is true.
	// staleRegions is not used when AlwaysReadPixelsFromGPU() returns true.
	staleRegions []image.Rectangle

	// pixelsCache is cached byte slices for pixels.
	// pixelsCache is just a cache to avoid allocations (#2375).
	//
	// A key is the region and a value is a byte slice for the region.
	//
	// It is fine to reuse the same byte slice for the same region for basePixels,
	// as old pixels for the same region will be invalidated at basePixel.AddOrReplace.
	pixelsCache map[image.Rectangle][]byte

	// regionsCache is cached regions.
	// regionsCache is just a cache to avoid allocations (#2375).
	regionsCache []image.Rectangle

	imageType ImageType
}

// NewImage creates an emtpy image with the given size.
//
// The returned image is cleared.
//
// Note that Dispose is not called automatically.
func NewImage(width, height int, imageType ImageType) *Image {
	if !graphicsDriverInitialized {
		panic("restorable: graphics driver must be ready at NewImage but not")
	}

	i := &Image{
		image:     graphicscommand.NewImage(width, height, imageType == ImageTypeScreen),
		width:     width,
		height:    height,
		imageType: imageType,
	}

	// This needs to use 'InternalSize' to render the whole region, or edges are unexpectedly cleared on some
	// devices.
	iw, ih := i.image.InternalSize()
	clearImage(i.image, image.Rect(0, 0, iw, ih))
	theImages.add(i)
	return i
}

// Extend extends the image by the given size.
// Extend creates a new image with the given size and copies the pixels of the given source image.
// Extend disposes itself after its call.
func (i *Image) Extend(width, height int) *Image {
	if i.width >= width && i.height >= height {
		return i
	}

	newImg := NewImage(width, height, i.imageType)

	// Use DrawTriangles instead of WritePixels because the image i might be stale and not have its pixels
	// information.
	srcs := [graphics.ShaderSrcImageCount]*Image{i}
	sw, sh := i.image.InternalSize()
	vs := make([]float32, 4*graphics.VertexFloatCount)
	graphics.QuadVerticesFromDstAndSrc(vs, 0, 0, float32(sw), float32(sh), 0, 0, float32(sw), float32(sh), 1, 1, 1, 1)
	is := graphics.QuadIndices()
	dr := image.Rect(0, 0, sw, sh)
	newImg.DrawTriangles(srcs, vs, is, graphicsdriver.BlendCopy, dr, [graphics.ShaderSrcImageCount]image.Rectangle{}, NearestFilterShader, nil, graphicsdriver.FillRuleFillAll)
	i.Dispose()

	return newImg
}

func clearImage(i *graphicscommand.Image, region image.Rectangle) {
	vs := make([]float32, 4*graphics.VertexFloatCount)
	graphics.QuadVerticesFromDstAndSrc(vs, float32(region.Min.X), float32(region.Min.Y), float32(region.Max.X), float32(region.Max.Y), 0, 0, 0, 0, 0, 0, 0, 0)
	is := graphics.QuadIndices()
	i.DrawTriangles([graphics.ShaderSrcImageCount]*graphicscommand.Image{}, vs, is, graphicsdriver.BlendClear, region, [graphics.ShaderSrcImageCount]image.Rectangle{}, clearShader.shader, nil, graphicsdriver.FillRuleFillAll)
}

// makeStale makes the image stale.
func (i *Image) makeStale(rect image.Rectangle) {
	i.stale = true
}

// ClearPixels clears the specified region by WritePixels.
func (i *Image) ClearPixels(region image.Rectangle) {
	i.WritePixels(nil, region)
}

func (i *Image) needsRestoring() bool {
	return i.imageType == ImageTypeRegular
}

// WritePixels replaces the image pixels with the given pixels slice.
//
// The specified region must not be overlapped with other regions by WritePixels.
func (i *Image) WritePixels(pixels *graphics.ManagedBytes, region image.Rectangle) {
	if region.Dx() <= 0 || region.Dy() <= 0 {
		panic("restorable: width/height must be positive")
	}
	w, h := i.width, i.height
	if !region.In(image.Rect(0, 0, w, h)) {
		panic(fmt.Sprintf("restorable: out of range %v", region))
	}

	// TODO: Avoid making other images stale if possible. (#514)
	// For this purpose, images should remember which part of that is used for DrawTriangles.
	theImages.makeStaleIfDependingOn(i)

	if pixels != nil {
		i.image.WritePixels(pixels, region)
	} else {
		clearImage(i.image, region)
	}

	// Even if the image is already stale, call makeStale to extend the stale region.
	i.makeStale(region)
}

// DrawTriangles draws triangles with the given image.
//
// The vertex floats are:
//
//	0: Destination X in pixels
//	1: Destination Y in pixels
//	2: Source X in texels
//	3: Source Y in texels
//	4: Color R [0.0-1.0]
//	5: Color G
//	6: Color B
//	7: Color Y
func (i *Image) DrawTriangles(srcs [graphics.ShaderSrcImageCount]*Image, vertices []float32, indices []uint32, blend graphicsdriver.Blend, dstRegion image.Rectangle, srcRegions [graphics.ShaderSrcImageCount]image.Rectangle, shader *Shader, uniforms []uint32, fillRule graphicsdriver.FillRule) {
	if len(vertices) == 0 {
		return
	}
	theImages.makeStaleIfDependingOn(i)

	// Even if the image is already stale, call makeStale to extend the stale region.
	i.makeStale(dstRegion)

	var imgs [graphics.ShaderSrcImageCount]*graphicscommand.Image
	for i, src := range srcs {
		if src == nil {
			continue
		}
		imgs[i] = src.image
	}
	i.image.DrawTriangles(imgs, vertices, indices, blend, dstRegion, srcRegions, shader.shader, uniforms, fillRule)
}

func (i *Image) ReadPixels(graphicsDriver graphicsdriver.Graphics, pixels []byte, region image.Rectangle) error {
	if err := i.image.ReadPixels(graphicsDriver, []graphicsdriver.PixelsArgs{
		{
			Pixels: pixels,
			Region: region,
		},
	}); err != nil {
		return err
	}
	return nil
}

// makeStaleIfDependingOn makes the image stale if the image depends on target.
func (i *Image) makeStaleIfDependingOn(target *Image) {
	if i.stale {
		return
	}
	if i.dependsOn(target) {
		// There is no new region to make stale.
		i.makeStale(image.Rectangle{})
	}
}

// makeStaleIfDependingOnShader makes the image stale if the image depends on shader.
func (i *Image) makeStaleIfDependingOnShader(shader *Shader) {
	if i.stale {
		return
	}
	if i.dependsOnShader(shader) {
		// There is no new region to make stale.
		i.makeStale(image.Rectangle{})
	}
}

// dependsOn reports whether the image depends on target.
func (i *Image) dependsOn(target *Image) bool {
	return false
}

// dependsOnShader reports whether the image depends on shader.
func (i *Image) dependsOnShader(shader *Shader) bool {
	return false
}

// dependingImages returns all images that is depended on the image.
func (i *Image) dependingImages() map[*Image]struct{} {
	r := map[*Image]struct{}{}
	return r
}

// hasDependency returns a boolean value indicating whether the image depends on another image.
func (i *Image) hasDependency() bool {
	return false
}

// Restore restores *graphicscommand.Image from the pixels using its state.
func (i *Image) restore(graphicsDriver graphicsdriver.Graphics) error {
	w, h := i.width, i.height
	// Do not dispose the image here. The image should be already disposed.

	switch i.imageType {
	case ImageTypeScreen:
		// The screen image should also be recreated because framebuffer might
		// be changed.
		i.image = graphicscommand.NewImage(w, h, true)
		i.basePixels.Dispose()
		i.basePixels = Pixels{}
		i.stale = false
		i.staleRegions = i.staleRegions[:0]
		return nil
	case ImageTypeVolatile:
		i.image = graphicscommand.NewImage(w, h, false)
		iw, ih := i.image.InternalSize()
		clearImage(i.image, image.Rect(0, 0, iw, ih))
		return nil
	}

	if i.stale {
		panic("restorable: pixels must not be stale when restoring")
	}

	gimg := graphicscommand.NewImage(w, h, false)
	// Clear the image explicitly.
	iw, ih := gimg.InternalSize()
	clearImage(gimg, image.Rect(0, 0, iw, ih))

	i.basePixels.Apply(gimg)

	i.image = gimg
	i.stale = false
	i.staleRegions = i.staleRegions[:0]
	return nil
}

// Dispose disposes the image.
//
// After disposing, calling the function of the image causes unexpected results.
func (i *Image) Dispose() {
	theImages.remove(i)
	i.image.Dispose()
	i.image = nil
	i.basePixels.Dispose()
	i.basePixels = Pixels{}
	i.pixelsCache = nil
	i.stale = false
	i.staleRegions = i.staleRegions[:0]
}

func (i *Image) Dump(graphicsDriver graphicsdriver.Graphics, path string, blackbg bool, rect image.Rectangle) (string, error) {
	return i.image.Dump(graphicsDriver, path, blackbg, rect)
}

func (i *Image) InternalSize() (int, int) {
	return i.image.InternalSize()
}

// appendRegionRemovingDuplicates adds a region to a given list of regions,
// but removes any duplicate between the newly added region and any existing regions.
//
// In case the newly added region is fully contained in any pre-existing region, this function does nothing.
// Otherwise, any pre-existing regions that are fully contained in the newly added region are removed.
//
// This is done to avoid unnecessary reading pixels from GPU.
func appendRegionRemovingDuplicates(regions *[]image.Rectangle, region image.Rectangle) {
	for _, r := range *regions {
		if region.In(r) {
			// The newly added rectangle is fully contained in one of the input regions.
			// Nothing to add.
			return
		}
	}
	// Separate loop, as regions must not get mutated before above return.
	n := 0
	for _, r := range *regions {
		if r.In(region) {
			continue
		}
		(*regions)[n] = r
		n++
	}
	*regions = append((*regions)[:n], region)
}
