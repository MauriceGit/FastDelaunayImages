package main

import (
	"fmt"
	"math"
	"math/rand"

	sc "mtSweepCircle"
	//sc "github.com/MauriceGit/sweepcircle"

	//v "mtVector"
	"os"
	"time"

	fgm "github.com/fogleman/delaunay"
	"github.com/fogleman/gg"

	//"strconv"
	//"runtime/pprof"
	//"flag"
	//"log"
	//"runtime"
	//"image"
	//"image/png"

	"github.com/pkg/profile"
	//"image/draw"
)

func drawImageVoronoi(d sc.Voronoi, imageName string, drawDetails bool) {
	var scale float64 = 1.0
	var imageSizeY float64 = 2000 * 3
	var imageSizeX float64 = 2000 * 3
	dc := gg.NewContext(int(imageSizeX), int(imageSizeY))

	// Background filling in white
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.Scale(2.0*3, 2.0*3)

	dc.SetLineWidth(0.6)
	//for i, e := range d.Edges {
	for i := 0; i < len(d.Edges); i += 2 {
		e := d.Edges[i]

		if e == sc.EmptyE {
			continue
		}

		dc.SetRGB(0, 0, 0)

		var v1 sc.Vector
		var v2 sc.Vector
		if e.VOrigin != sc.EmptyVertex {
			v1 = d.Vertices[e.VOrigin].Pos
		} else {
			dc.SetRGB(1, 0, 0)
			v1 = sc.Add(d.Vertices[d.Edges[e.ETwin].VOrigin].Pos, d.Edges[i].TmpEdge.Dir)
			v1 = sc.Add(d.Edges[i].TmpEdge.Pos, d.Edges[i].TmpEdge.Dir)
		}
		if d.Edges[e.ETwin].VOrigin != sc.EmptyVertex {
			v2 = d.Vertices[d.Edges[e.ETwin].VOrigin].Pos
		} else {
			dc.SetRGB(0, 1, 0)
			v2 = sc.Add(d.Edges[i].TmpEdge.Pos, d.Edges[i].TmpEdge.Dir)
		}

		dc.DrawLine(v1.X*scale, v1.Y*scale, v2.X*scale, v2.Y*scale)
		dc.Stroke()
		//dc.SetRGB(0, 0, 0)

		if drawDetails {
			dc.SetRGB(0, 0.5, 0)
			middleP := sc.Vector{(v1.X + v2.X) / 2., (v1.Y + v2.Y) / 2.}

			crossP := sc.Perpendicular(sc.Sub(v1, v2))

			crossP.Div(sc.Length(crossP))
			crossP.Mult(15.)

			middleP.Add(crossP)

			i = i
			s := fmt.Sprintf("(%d)", i)
			dc.DrawStringAnchored(s, middleP.X, middleP.Y, 0.5, 0.5)
		}
	}
	dc.Stroke()

	//dc.SetLineWidth(1.0)
	dc.SetRGB(0.0, 0.1, 0.2)
	for i, v := range d.Vertices {

		if v == sc.EmptyV {
			continue
		}

		//dc.DrawCircle(v.Pos.X*scale, v.Pos.Y*scale, 2)
		//dc.DrawPoint(v.Pos.X*scale, v.Pos.Y*scale, 2.8)
		//dc.Fill()

		if drawDetails {
			s := fmt.Sprintf("(%d)", i)
			dc.DrawStringAnchored(s, v.Pos.X-10, v.Pos.Y-10, 0.5, 0.5)
		}
	}
	dc.Fill()

	if drawDetails {
		dc.SetRGB(0.8, 0.0, 0.0)
		for i, f := range d.Faces {

			if f == sc.EmptyF {
				continue
			}

			//dc.DrawCircle(v.Pos.X*scale, v.Pos.Y*scale, 2)
			//dc.DrawPoint(v.Pos.X*scale, v.Pos.Y*scale, 2.8)
			//dc.Fill()

			i = i
			s := fmt.Sprintf("%d", i)

			v0 := d.Vertices[d.Edges[f.EEdge].VOrigin].Pos
			v1 := d.Vertices[d.Edges[d.Edges[f.EEdge].ENext].VOrigin].Pos
			v2 := d.Vertices[d.Edges[d.Edges[d.Edges[f.EEdge].ENext].ENext].VOrigin].Pos

			center := sc.Add(v0, v1)
			center = sc.Add(center, v2)
			center = sc.Mult(center, 1./3.)
			s = s

			//			dc.DrawStringAnchored(s, center.X, center.Y, 0.5, 0.5)
		}
	}

	//dc.SetRGB(1, 1, 0)
	//dc.DrawCircle(432, imageSizeY-894, 5)
	//dc.DrawCircle(599, imageSizeY-532, 5)
	//dc.DrawCircle(501, imageSizeY-578, 5)
	//dc.Fill()

	dc.SavePNG(imageName + ".png")
}

func drawImage(d sc.Delaunay, imageName string, drawDetails bool) {
	var scale float64 = 1.0
	var imageSizeY float64 = 2000
	var imageSizeX float64 = 2000
	dc := gg.NewContext(int(imageSizeX), int(imageSizeY))

	// Background filling in white
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.Scale(2.0, 2.0)

	dc.SetLineWidth(0.8)

	// Draws 10 horizontal and vertical lines over the whole image and labels them on the side.
	for i := 1; i < 10; i++ {

		//		x := float64(i) * 100 * scale
		//		y := float64(i) * 100 * scale

		//		dc.SetRGB(1, 0.5, 0.5)
		//		dc.DrawLine(0, y, imageSizeX, y)
		//		dc.Stroke()
		//		dc.DrawLine(x, 0, x, imageSizeY)
		//		dc.Stroke()

		//		dc.SetRGB(1, 0.0, 0.0)
		//		// X axis
		//		//dc.DrawString(strconv.Itoa(int(x)), x+10, (imageSizeY/2)-10)
		//		// Y axis
		//		//dc.DrawString(strconv.Itoa(int(imageSizeY-y)), 10, y-10)

	}

	dc.SetLineWidth(0.6)
	for i, e := range d.Edges {

		if e == sc.EmptyE {
			continue
		}

		dc.SetRGB(0, 0, 0)

		v1 := d.Vertices[e.VOrigin].Pos
		v2 := d.Vertices[d.Edges[e.ETwin].VOrigin].Pos

		dc.DrawLine(v1.X*scale, v1.Y*scale, v2.X*scale, v2.Y*scale)

		if drawDetails {
			dc.SetRGB(0, 0.5, 0)
			middleP := sc.Vector{(v1.X + v2.X) / 2., (v1.Y + v2.Y) / 2.}

			crossP := sc.Perpendicular(sc.Sub(v1, v2))

			crossP.Div(sc.Length(crossP))
			crossP.Mult(15.)

			middleP.Add(crossP)

			i = i
			s := fmt.Sprintf("(%d)", i)
			dc.DrawStringAnchored(s, middleP.X, middleP.Y, 0.5, 0.5)
		}
	}
	dc.Stroke()

	//dc.SetLineWidth(1.0)
	dc.SetRGB(0.0, 0.1, 0.2)
	for i, v := range d.Vertices {

		if v == sc.EmptyV {
			continue
		}

		//dc.DrawCircle(v.Pos.X*scale, v.Pos.Y*scale, 2)
		dc.DrawPoint(v.Pos.X*scale, v.Pos.Y*scale, 2.8)
		//dc.Fill()

		if drawDetails {
			s := fmt.Sprintf("(%d)", i)
			dc.DrawStringAnchored(s, v.Pos.X-10, v.Pos.Y-10, 0.5, 0.5)
		}
	}
	dc.Fill()

	if drawDetails {
		dc.SetRGB(0.8, 0.0, 0.0)
		for i, f := range d.Faces {

			if f == sc.EmptyF {
				continue
			}

			//dc.DrawCircle(v.Pos.X*scale, v.Pos.Y*scale, 2)
			//dc.DrawPoint(v.Pos.X*scale, v.Pos.Y*scale, 2.8)
			//dc.Fill()

			i = i
			s := fmt.Sprintf("%d", i)

			v0 := d.Vertices[d.Edges[f.EEdge].VOrigin].Pos
			v1 := d.Vertices[d.Edges[d.Edges[f.EEdge].ENext].VOrigin].Pos
			v2 := d.Vertices[d.Edges[d.Edges[d.Edges[f.EEdge].ENext].ENext].VOrigin].Pos

			center := sc.Add(v0, v1)
			center = sc.Add(center, v2)
			center = sc.Mult(center, 1./3.)
			s = s

			//			dc.DrawStringAnchored(s, center.X, center.Y, 0.5, 0.5)
		}
	}

	//dc.SetRGB(1, 1, 0)
	//dc.DrawCircle(432, imageSizeY-894, 5)
	//dc.DrawCircle(599, imageSizeY-532, 5)
	//dc.DrawCircle(501, imageSizeY-578, 5)
	//dc.Fill()

	dc.SavePNG(imageName + ".png")
}

func drawFgmImage(points []fgm.Point, triangulation *fgm.Triangulation, imageName string) {

	W := 2048
	H := 2048

	nextHalfEdge := func(e int) int {
		if e%3 == 2 {
			return e - 2
		}
		return e + 1
	}

	// compute point bounds for rendering
	min := points[0]
	max := points[0]
	for _, p := range points {
		min.X = math.Min(min.X, p.X)
		min.Y = math.Min(min.Y, p.Y)
		max.X = math.Max(max.X, p.X)
		max.Y = math.Max(max.Y, p.Y)
	}

	size := fgm.Point{max.X - min.X, max.Y - min.Y}
	center := fgm.Point{min.X + size.X/2, min.Y + size.Y/2}
	scale := math.Min(float64(W)/size.X, float64(H)/size.Y) * 0.9

	// render points and edges
	dc := gg.NewContext(W, H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	dc.Translate(float64(W/2), float64(H/2))
	dc.Scale(scale, scale)
	dc.Translate(-center.X, -center.Y)

	ts := triangulation.Triangles
	hs := triangulation.Halfedges
	for i, h := range hs {
		if i > h {
			p := points[ts[i]]
			q := points[ts[nextHalfEdge(i)]]
			dc.DrawLine(p.X, p.Y, q.X, q.Y)
		}
	}
	dc.Stroke()

	for _, p := range points {
		dc.DrawPoint(p.X, p.Y, 5)
	}
	dc.Fill()

	for _, p := range triangulation.ConvexHull {
		dc.LineTo(p.X, p.Y)
	}
	dc.ClosePath()
	dc.SetLineWidth(5)
	dc.Stroke()

	dc.SavePNG(imageName + ".png")
}

func triangulate(myPoints []sc.Vector, fgmPoints []fgm.Point, renderImage, profileMT, profileFgm bool, imageName string) {

	var profiling interface {
		Stop()
	}

	if profileFgm {
		profiling = profile.Start(profile.CPUProfile)
	}

	start := time.Now()
	var triangulationFgm *fgm.Triangulation
	var err error
	if fgmPoints != nil && len(fgmPoints) > 0 {
		triangulationFgm, err = fgm.Triangulate(fgmPoints)
		if err != nil {
			fmt.Printf("Fogleman encountered a serious error: %v\n", err)
		}
	}
	binTimeFgm := time.Since(start).Nanoseconds()

	if !profileFgm && profileMT {
		profiling = profile.Start(profile.CPUProfile)
	}

	start = time.Now()
	var triangulationMT sc.Delaunay
	if myPoints != nil && len(myPoints) > 0 {
		triangulationMT = sc.Triangulate(myPoints)
	}
	binTime := time.Since(start).Nanoseconds()

	start = time.Now()
	var voronoiMT sc.Voronoi
	if myPoints != nil && len(myPoints) > 0 {
		voronoiMT = triangulationMT.CreateVoronoi()
	}
	binTimeVoronoi := time.Since(start).Nanoseconds()

	if profileFgm || profileMT {
		profiling.Stop()
	}

	errFgm := triangulationFgm.Validate()
	errMT := triangulationMT.Verify()
	errMTV := voronoiMT.Verify()
	if errFgm != nil {
		//fmt.Printf("Fogleman encountered an error: %v\n", errFgm)
	}
	if errMT != nil {
		fmt.Printf("SweepCircle encountered an error: %v\n", errMT)
	}
	if errMTV != nil {
		fmt.Printf("Voronoi encountered an error: %v\n", errMTV)
	}

	fmt.Printf("Triangulation (ms): %.8f, Voronoi (ms): %.8f, Fogleman (ms): %.8f\n", float64(binTime)/1000000.0, float64(binTimeVoronoi)/1000000.0, float64(binTimeFgm)/1000000.0)

	if renderImage {
		drawImage(triangulationMT, imageName+"_mt_", false)
		drawImageVoronoi(voronoiMT, imageName+"_mt_voronoi_", false)
		//drawFgmImage(fgmPoints, triangulationFgm, imageName+"_fgm_")
	}

}

func toFgmList(points []sc.Vector) []fgm.Point {
	newPoints := make([]fgm.Point, len(points), len(points))
	for i, v := range points {
		newPoints[i] = fgm.Point{v.X, v.Y}
	}
	return newPoints
}

func testRandom(count int) ([]sc.Vector, []fgm.Point) {
	fmt.Printf("===========================\n")
	fmt.Printf("=== test random\n")
	fmt.Printf("===========================\n")

	scale := 1000.0
	margin := 10.0

	var seed int64 = time.Now().UTC().UnixNano()
	seed = seed
	fmt.Fprintf(os.Stderr, "Seed: %v\n", seed)
	r := rand.New(rand.NewSource(seed))
	var pointList []sc.Vector

	for i := 0; i < count; i++ {
		pointList = append(pointList, sc.Vector{r.Float64()*(scale-2*margin) + margin, r.Float64()*(scale-2*margin) + margin})
	}

	return pointList, toFgmList(pointList)
}

func testCircle(count int) ([]sc.Vector, []fgm.Point) {
	fmt.Printf("===========================\n")
	fmt.Printf("=== test circle\n")
	fmt.Printf("===========================\n")

	var pointList []sc.Vector
	center := sc.Vector{500, 500}
	radius := 450.

	for i := 0; i < count; i++ {
		newI := sc.DegToRad(float64(i) / float64(count) * 360.)
		v := sc.Vector{center.X + radius*math.Cos(float64(newI)), center.Y + radius*math.Sin(float64(newI))}

		pointList = append(pointList, v)
	}

	return pointList, toFgmList(pointList)
}

func testDoubleCircle(count int) ([]sc.Vector, []fgm.Point) {
	fmt.Printf("===========================\n")
	fmt.Printf("=== test double circle\n")
	fmt.Printf("===========================\n")

	var pointList []sc.Vector
	center := sc.Vector{500, 500}
	radius := 350.

	for i := 0; i < count/2; i++ {
		newI := sc.DegToRad(float64(i) / float64(count/2) * float64(i))
		v := sc.Vector{center.X + radius*math.Cos(float64(newI)), center.Y + radius*math.Sin(float64(newI))}

		pointList = append(pointList, v)
	}

	radius = 450.

	for i := 0; i < count/2; i++ {
		newI := sc.DegToRad(float64(i) / float64(count/2) * 360.)
		v := sc.Vector{center.X + radius*math.Cos(float64(newI)), center.Y + radius*math.Sin(float64(newI))}
		pointList = append(pointList, v)
	}

	//pointList = append(pointList, sc.Vector{500, 500})

	return pointList, toFgmList(pointList)
}

func testWaveCenterMirrored(count int) ([]sc.Vector, []fgm.Point) {
	fmt.Printf("===========================\n")
	fmt.Printf("=== test wave center mirrored\n")
	fmt.Printf("===========================\n")

	var pointList []sc.Vector

	for i := 0; i <= count; i++ {

		newI := sc.DegToRad(float64(i) / float64(count) * 360. * 5.0)

		v := sc.Vector{float64(i)/float64(count)*900 + 50, math.Sin(newI)*450. + 500}
		//fmt.Println(v)
		pointList = append(pointList, v)
	}

	return pointList, toFgmList(pointList)
}

func testWave(count int) ([]sc.Vector, []fgm.Point) {
	fmt.Printf("===========================\n")
	fmt.Printf("=== test wave\n")
	fmt.Printf("===========================\n")

	var pointList []sc.Vector

	for i := 0; i <= count; i++ {
		newI := sc.DegToRad(float64(i) / float64(count) * 360.)

		v := sc.Vector{float64(i)/float64(count)*900 + 50, math.Sin(newI)*450. + 500}
		//fmt.Println(v)
		pointList = append(pointList, v)
	}

	return pointList, toFgmList(pointList)
}

func testTiltedGrid(count int, tiltAngle float64) ([]sc.Vector, []fgm.Point) {
	fmt.Printf("===========================\n")
	fmt.Printf("=== test tilted grid %.2f\n", tiltAngle)
	fmt.Printf("===========================\n")

	angle := sc.DegToRad(tiltAngle)

	var pointList []sc.Vector

	for x := 0; x <= count; x++ {
		newX := float64(x)/float64(count)*900. + 50.
		for y := 0; y <= count; y++ {
			newY := float64(y)/float64(count)*900. + 50.
			tiltedX := (newX-500.)*math.Cos(angle) - (newY-500.)*math.Sin(angle)
			tiltedY := (newY-500.)*math.Cos(angle) + (newX-500.)*math.Sin(angle)
			tiltedX += 500.
			tiltedY += 500.
			v := sc.Vector{tiltedX, tiltedY}
			pointList = append(pointList, v)
		}
	}

	return pointList, toFgmList(pointList)
}

func testPoisson(count int) ([]sc.Vector, []fgm.Point) {
	fmt.Printf("===========================\n")
	fmt.Printf("=== test random\n")
	fmt.Printf("===========================\n")

	scale := 1000.0
	margin := 10.0

	var seed int64 = time.Now().UTC().UnixNano()
	seed = seed
	fmt.Fprintf(os.Stderr, "Seed: %v\n", seed)
	pointList := CreateFastPoissonDiscPoints(int(float64(count)*1.8), scale, scale, margin, 30, seed)

	fmt.Printf("poisson list count: %v\n", len(pointList))

	return pointList, toFgmList(pointList)
}

func main() {

	var myP []sc.Vector
	var fgmP []fgm.Point
	myP, fgmP = testRandom(1000)
	triangulate(myP, fgmP, false, false, false, "random")

	myP, fgmP = testPoisson(100)
	triangulate(myP, fgmP, true, false, false, "poisson")

	/*
		myP, fgmP = testTiltedGrid(500, 0.0)
		triangulate(myP, fgmP, false, false, false, "tilted_0")

		myP, fgmP = testTiltedGrid(500, 89.0)
		triangulate(myP, fgmP, false, false, false, "tilted_89")
		myP, fgmP = testTiltedGrid(500, 45.0)
		triangulate(myP, fgmP, false, false, false, "tilted_45")
		myP, fgmP = testCircle(5000)
		triangulate(myP, fgmP, false, false, false, "circle")
		myP, fgmP = testDoubleCircle(10000)
		triangulate(myP, fgmP, false, false, false, "double_circle")
		myP, fgmP = testWaveCenterMirrored(10000)
		triangulate(myP, fgmP, false, false, false, "wave_mirrored")
		myP, fgmP = testWave(10000)
		triangulate(myP, fgmP, false, false, false, "wave")
	*/

	///=========== Frontier: Slice ==================================///

	// Correct timing and Interpolation Search. Seconds.
	// 10       points  0.00005763
	// 100      points  0.00097302
	// 1000     points  0.01580784
	// 10000    points  0.19852256
	// 100000   points  2.56392177
	// 1000000  points  31.5917846

	// Implemented a hardcoded determinant for maximum speed!
	// 10       points  0.00003697
	// 100      points  0.00026984
	// 1000     points  0.00297420
	// 10000    points  0.03873088
	// 100000   points  0.65939308
	// 1000000  points  5.31894416
	// 2000000  points  11.72728091
	// 5000000  points  33.50491853

	///=========== Frontier: 2,3-Tree ===============================///

	// Changed Frontier structure to first 2,3-tree implementation
	// ===========================
	// === test_unknown_problem_random
	// ===========================
	// count: 1000000
	// Seed:  1527368018753283347
	// Triangulation in Milliseconds: 16.10789566

	// Optimized the 2,3-tree some more (not exactly sure, but before introducing memory recycling and its own memory manager)
	// ===========================
	// === test_unknown_problem_random
	// ===========================
	// count: 1000000
	// Seed:  1527368018753283347
	// Triangulation in Milliseconds: 11.60121581

	// Optimized 2,3-tree with node recycling and memory management.
	// 10       points  0.08035200
	// 100      points  0.00056354
	// 1000     points  0.00479936
	// 10000    points  0.05738459
	// 100000   points  0.42946737
	// 1000000  points  7.71770482
	// 2000000  points  16.2359380

	// Some finetuning.
	// 10       points  0.00008035
	// 100      points  0.00047035
	// 1000     points  0.00444172
	// 10000    points  0.05291956
	// 100000   points  0.57214100
	// 1000000  points  6.40743041
	// 2000000  points  13.4764279
	// 5000000  points  36.0792809

	///=========== Frontier: External SkipList ======================///

	// With extern SkipList and lots of .Seek()
	// 10       points  0.00010022
	// 100      points  0.00092334
	// 1000     points  0.01044300
	// 10000    points  0.09309602
	// 100000   points  0.99515928
	// 1000000  points  11.8651977
	// 2000000  points  25.0317357
	// 5000000  points  64.6695457

	///=========== Frontier: My own SkipList ======================///

	// 10       points  0.00008780400
	// 100      points  0.00039038600
	// 1000     points  0.00341958600
	// 10000    points  0.03582979400
	// 100000   points  0.39232895800
	// 1000000  points  4.24569866200
	// 2000000  points  9.37676779900
	// 5000000  points  24.3403684720

	///=========== Maurice Laptop -- with Profiling ======================///

	///=========== Github-Stand ======================///
	// 1000000  points  9.17114756000
	///=========== 2D-Points ======================///
	// 1000000  points  8.27235414300
	///=========== Length squared ======================///
	// 1000000  points  8.11650684100
	///=========== Angle rad ======================///
	// 1000000  points  8.05951195700
	///=========== Fast acos ======================///
	// 1000000  points  7.65465191200
	///=========== Skiplist height adjustment ======================///
	// 1000000  points  7.45683164800
	///=========== Triangle valid from Fogleman ======================///
	// 1000000  points  6.33755930000
	///=========== Multithreaded preparation ======================///
	// 1000000  points  6.17650878400

	///=========== New PC -- No profiling ======================///
	// ArrayMap instead of skiplist
	// 2D point math
	// Fast Acos approximation
	// Only use rad angle
	// Squared distance
	// Triangle validation from Fogleman
	///=============================================

	// 10       points  0.00014780400   0.00003631600
	// 100      points  0.00043934400   0.00019591000
	// 1000     points  0.00326873000   0.00211053000
	// 10000    points  0.00757377300   0.00607904100
	// 100000   points  0.07343298500   0.06884289200
	// 1000000  points  0.80409688000   0.96240224000
	// 2000000  points  1.67112829300   2.17478010300
	// 5000000  points  4.37922068100   6.26767477000
	// 10000000 points  9.12704656300   14.4532323540

	///=============================================
	// Fast arrayMap (with mem pool)
	// Custom Quicksort
	// No Delaunay edge re-creation
	///=============================================

	// 10       points  0.00019367300   0.00003631600
	// 100      points  0.00039306300   0.00019591000
	// 1000     points  0.00283277200   0.00211053000
	// 10000    points  0.00602252900   0.00607904100
	// 100000   points  0.05958530500   0.06884289200
	// 1000000  points  0.64559663300   0.96240224000
	// 2000000  points  1.33272351500   2.17478010300
	// 5000000  points  3.55517294500   6.26767477000
	// 10000000 points  7.36102560900   14.4532323540

}
