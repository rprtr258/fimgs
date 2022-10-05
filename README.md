# fimgs - image filters tool

## Install
```bash
go install github.com/rprtr258/fimgs/cmd/fimgs@latest
```

## Usage
```
NAME:
   fimgs - Applies filter to image.

USAGE:
   Applies filter to image and saves new image.

COMMANDS:
   sharpen          Sharpen filter.
   edgeenhance      Edgeenhance filter.
   edgedetect1      Edgedetect1 filter.
   verticallines    Verticallines filter.
   horizontallines  Horizontallines filter.
   blur             Blur filter.
   weakblur         Weakblur filter.
   emboss           Emboss filter.
   edgedetect2      Edgedetect2 filter.
   cluster          Cluster colors.
   quadtree         Quad tree filter.
   shader           Shader filter.
   hilbert          Hilbert curve filter.
   hilbertdarken    Hilbert darken curve filter.
   zcurve           Z curve filter.
   median           Median filter.
   help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h               show help (default: false)
   --image value, -i value  input image filename
   

```

## Examples
||||
|-|-|-|
|![](img/static/orig.png)|![](img/static/blur.png)|![](img/static/cluster.png)|
|orig|blur|cluster|
|![](img/static/edgedetect1.png)|![](img/static/edgedetect2.png)|![](img/static/edgeenhance.png)|
|edgedetect1|edgedetect2|edgeenhance|
|![](img/static/emboss.png)|![](img/static/hilbert.png)|![](img/static/hilbertdarken.png)|
|emboss|[hilbert](https://habr.com/en/post/135344/)|hilbertdarken|
|![](img/static/horizontallines.png)|![](img/static/median.png)|![](img/static/quadtree.png)|
|horizontallines|[median](https://en.wikipedia.org/wiki/Median_filter)|[quadtree](https://habr.com/en/post/280674/)|
|![](img/static/shader_rgb.png)|![](img/static/sharpen.png)|![](img/static/verticallines.png)|
|shader/rgb|sharpen|verticallines|
|![](img/static/weakblur.png)|![](img/static/zcurve.png)||
|weakblur|zcurve||

