# fimgs - image filters tool

## Install

```bash
go install github.com/rprtr258/fimgs/cmd/fimgs@latest
```

## Usage

```php
NAME:
   fimgs - Applies filter to image

USAGE:
   Applies filter to image and saves new image

COMMANDS:
   blur             Blur filter
   cluster          Cluster colors
   edgedetect1      Edgedetect1 filter
   edgedetect2      Edgedetect2 filter
   edgeenhance      Edgeenhance filter
   emboss           Emboss filter
   hilbert          Hilbert curve filter
   hilbertdarken    Hilbert darken curve filter
   horizontallines  Horizontallines filter
   median           Median filter
   quadtree         Quad tree filter
   shader           Shader filter
   sharpen          Sharpen filter
   verticallines    Verticallines filter
   weakblur         Weakblur filter
   zcurve           Z curve filter

GLOBAL OPTIONS:
   --image value, -i value  input image filename
   --help, -h               show help

```

## Examples

||||
|-|-|-|
|![](./img/static/blur.png)|![](./img/static/cluster.png)|![](./img/static/edgedetect1.png)|
|blur.png|cluster.png|edgedetect1.png|
|![](./img/static/edgedetect2.png)|![](./img/static/edgeenhance.png)|![](./img/static/emboss.png)|
|edgedetect2.png|edgeenhance.png|emboss.png|
|![](./img/static/hilbert.png)|![](./img/static/hilbertdarken.png)|![](./img/static/horizontallines.png)|
|hilbert.png|hilbertdarken.png|horizontallines.png|
|![](./img/static/median.png)|![](./img/static/orig.png)|![](./img/static/quadtree.png)|
|median.png|orig.png|quadtree.png|
|![](./img/static/shader_rgb.png)|![](./img/static/sharpen.png)|![](./img/static/verticallines.png)|
|shader_rgb.png|sharpen.png|verticallines.png|
|![](./img/static/weakblur.png)|![](./img/static/zcurve.png)|
|weakblur.png|zcurve.png|
