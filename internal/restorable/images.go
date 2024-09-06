// Copyright 2017 The Ebiten Authors
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
	"github.com/hajimehoshi/ebiten/v2/internal/debug"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicscommand"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver"
)

// images is a set of Image objects.
type images struct {
	images     map[*Image]struct{}
	shaders    map[*Shader]struct{}
	lastTarget *Image
}

// theImages represents the images for the current process.
var theImages = &images{
	images:  map[*Image]struct{}{},
	shaders: map[*Shader]struct{}{},
}

func SwapBuffers(graphicsDriver graphicsdriver.Graphics) error {
	if debug.IsDebug {
		debug.FrameLogf("Internal image sizes:\n")
		imgs := make([]*graphicscommand.Image, 0, len(theImages.images))
		for i := range theImages.images {
			imgs = append(imgs, i.image)
		}
		graphicscommand.LogImagesInfo(imgs)
	}
	if err := graphicscommand.FlushCommands(graphicsDriver, true); err != nil {
		return err
	}
	return nil
}

// DumpImages dumps all the current images to the specified directory.
//
// This is for testing usage.
func DumpImages(graphicsDriver graphicsdriver.Graphics, dir string) (string, error) {
	images := make([]*graphicscommand.Image, 0, len(theImages.images))
	for img := range theImages.images {
		images = append(images, img.image)
	}

	return graphicscommand.DumpImages(images, graphicsDriver, dir)
}

// add adds img to the images.
func (i *images) add(img *Image) {
	i.images[img] = struct{}{}
}

func (i *images) addShader(shader *Shader) {
	i.shaders[shader] = struct{}{}
}

// remove removes img from the images.
func (i *images) remove(img *Image) {
	i.makeStaleIfDependingOn(img)
	delete(i.images, img)
}

func (i *images) removeShader(shader *Shader) {
	i.makeStaleIfDependingOnShader(shader)
	delete(i.shaders, shader)
}

// makeStaleIfDependingOn makes all the images stale that depend on target.
//
// When target is modified, all images depending on target can't be restored with target.
// makeStaleIfDependingOn is called in such situation.
func (i *images) makeStaleIfDependingOn(target *Image) {
	if target == nil {
		panic("restorable: target must not be nil at makeStaleIfDependingOn")
	}
	if i.lastTarget == target {
		return
	}
	i.lastTarget = target
	for img := range i.images {
		img.makeStaleIfDependingOn(target)
	}
}

// makeStaleIfDependingOn makes all the images stale that depend on shader.
func (i *images) makeStaleIfDependingOnShader(shader *Shader) {
	if shader == nil {
		panic("restorable: shader must not be nil at makeStaleIfDependingOnShader")
	}
	for img := range i.images {
		img.makeStaleIfDependingOnShader(shader)
	}
}

var graphicsDriverInitialized bool

// InitializeGraphicsDriverState initializes the graphics driver state.
func InitializeGraphicsDriverState(graphicsDriver graphicsdriver.Graphics) error {
	graphicsDriverInitialized = true
	return graphicscommand.InitializeGraphicsDriverState(graphicsDriver)
}

// MaxImageSize returns the maximum size of an image.
func MaxImageSize(graphicsDriver graphicsdriver.Graphics) int {
	return graphicscommand.MaxImageSize(graphicsDriver)
}
