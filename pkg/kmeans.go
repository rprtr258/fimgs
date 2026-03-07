package fimgs

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"simd/archsimd"
)

func makeColorArray(len int) []archsimd.Int64x4 {
	return make([]archsimd.Int64x4, len)
}

func minkowskiiDist(a, b archsimd.Int64x4) int64 {
	d := a.Sub(b).Abs()
	t := d.GetHi().Add(d.GetLo())
	return t.GetElem(0) + t.GetElem(1)
}

func initClusterCenters(pixelColors []archsimd.Int64x4, clustersCount int) []archsimd.Int64x4 {
	clustersCenters := makeColorArray(clustersCount)
	clustersCenters[clustersCount-1] = pixelColors[rand.Intn(len(pixelColors))]
	minClusterDistance := make([]int64, len(pixelColors))
	minClusterDistanceSum := int64(0)
	for i, pixelColor := range pixelColors {
		minClusterDistance[i] = minkowskiiDist(pixelColor, clustersCenters[clustersCount-1])
		minClusterDistanceSum += minClusterDistance[i]
	}
	for k := 0; k < clustersCount-1; k++ {
		x := rand.Int63n(minClusterDistanceSum)
		for i, pixelColor := range pixelColors {
			x -= minClusterDistance[i]
			if x < 0 {
				clustersCenters[k] = pixelColor
				break
			}
		}
		if k == clustersCount-1 {
			break
		}
		for i, pixelColor := range pixelColors {
			if newDistance := minkowskiiDist(pixelColor, clustersCenters[0]); newDistance < minClusterDistance[i] {
				minClusterDistanceSum += newDistance - minClusterDistance[i]
				minClusterDistance[i] = newDistance
			}
		}
	}
	return clustersCenters
}

func kmeansIters(clustersCenters, pixelColors []archsimd.Int64x4, clustersCount int) {
	batchMaxSize := int(math.Sqrt(float64(len(pixelColors))))
	sumAndCount := make([]archsimd.Int64x4, clustersCount) // sum of Rs, Gs, Bs and count
	for range 300 {
		for i := range sumAndCount {
			sumAndCount[i] = archsimd.Int64x4{}
		}

		k := rand.Intn(batchMaxSize) + 1
		for i := k; i < len(pixelColors); i += k {
			pixelColor := pixelColors[i]
			minCluster := 0
			minDist := minkowskiiDist(pixelColor, clustersCenters[0])
			for k := 1; k < clustersCount; k++ {
				newDist := minkowskiiDist(pixelColor, clustersCenters[k])
				if newDist < minDist {
					minCluster = k
					minDist = newDist
				}
			}
			sumAndCount[minCluster] = sumAndCount[minCluster].Add(pixelColor)
		}
		movement := int64(0)
		for i := range clustersCount {
			count := sumAndCount[i].GetLo().GetElem(0)
			const S = 62
			divisor := archsimd.BroadcastInt64x4((int64(1) << S) / count)
			v := sumAndCount[i].Mul(divisor).ShiftAllRight(S)
			movement += minkowskiiDist(clustersCenters[i], v)
			clustersCenters[i] = v
		}
		if movement < 100 {
			break
		}
	}
}

func ApplyKMeans(im image.Image, clustersCount int) image.RGBA {
	imageWidth := im.Bounds().Dx()
	pixelColors := makeColorArray(imageWidth * im.Bounds().Dy())
	for j := 0; j < im.Bounds().Dy(); j++ {
		for i := 0; i < im.Bounds().Dx(); i++ {
			k := i + j*imageWidth
			r, g, b, _ := im.At(i, j).RGBA()
			pixelColors[k] = archsimd.LoadInt64x4(&[4]int64{int64(r), int64(g), int64(b), 1})
		}
	}
	rand.Seed(0)
	clustersCenters := initClusterCenters(pixelColors, clustersCount)
	// TODO: try to sample mini-batches (random subdatasets)
	kmeansIters(clustersCenters, pixelColors, clustersCount)
	filtered_im := image.NewRGBA(im.Bounds())
	for j := 0; j < im.Bounds().Dy(); j++ {
		for i := 0; i < im.Bounds().Dx(); i++ {
			pixel := pixelColors[i+j*imageWidth]
			minCluster := 0
			minDist := minkowskiiDist(pixel, clustersCenters[0])
			for k := 1; k < clustersCount; k++ {
				dist := minkowskiiDist(pixel, clustersCenters[k])
				if dist < minDist {
					minCluster = k
					minDist = dist
				}
			}
			x := clustersCenters[minCluster].And(archsimd.BroadcastInt64x4(0xFF00))
			filtered_im.Set(i, j, color.RGBA{
				uint8((x.GetHi().GetElem(0)) >> 8),
				uint8((x.GetHi().GetElem(1)) >> 8),
				uint8((x.GetLo().GetElem(0)) >> 8),
				255,
			})
		}
	}
	return *filtered_im
}

// TODO: filter init is also validation?
func ApplyKMeansFilter(sourceImageFilename string, resultImageFilename string, clustersCount int) error {
	if clustersCount < 2 {
		return fmt.Errorf("'n' must be at least 2, you gave n=%d", clustersCount)
	}

	im, err := LoadImageFile(sourceImageFilename)
	if err != nil {
		return fmt.Errorf("error occured while loading image: %q", err)
	}

	filtered_im := ApplyKMeans(im, clustersCount)

	if err := saveImage(filtered_im, resultImageFilename); err != nil {
		return err
	}

	return nil
}
