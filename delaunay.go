package main

import (
    "fmt"
    "os"
    sc "mtSweepCircle"
    v  "mtVector"
    "math"
    "math/rand"
    "time"
    "github.com/fogleman/gg"
    fgm "github.com/fogleman/delaunay"
    //"strconv"
    //"runtime/pprof"
    //"flag"
    //"log"
    //"runtime"
    "github.com/pkg/profile"
    "image/png"
    "image"
    //"image/draw"
)


func calcEdgeColorFromVertices(img image.Image, bounds image.Rectangle, v1, v2 v.Vector) (float64, float64, float64) {
    xRange := float64(bounds.Max.X-bounds.Min.X)
    yRange := float64(bounds.Max.Y-bounds.Min.Y)
    r1, g1, b1, _ := img.At(int(v1.X/1000.0*xRange), int(v1.Y/1000.0*yRange)).RGBA()
    r2, g2, b2, _ := img.At(int(v2.X/1000.0*xRange), int(v2.Y/1000.0*yRange)).RGBA()

    return float64(((r1+r2)/2)>>8)/255., float64(((g1+g2)/2)>>8)/255., float64(((b1+b2)/2.0)>>8)/255.
}

func calcFaceColor(img image.Image, bounds image.Rectangle, v1, v2, v3 v.Vector) (float64, float64, float64) {
    xRange := float64(bounds.Max.X-bounds.Min.X)
    yRange := float64(bounds.Max.Y-bounds.Min.Y)

    r1, g1, b1, _ := img.At(int(v1.X/1000.0*xRange), int(v1.Y/1000.0*yRange)).RGBA()
    r2, g2, b2, _ := img.At(int(v2.X/1000.0*xRange), int(v2.Y/1000.0*yRange)).RGBA()
    r3, g3, b3, _ := img.At(int(v3.X/1000.0*xRange), int(v3.Y/1000.0*yRange)).RGBA()
    center := v.Div(v.Add(v.Add(v1,v2), v3), 3)
    rc, gc, bc, _ := img.At(int(center.X/1000.0*xRange), int(center.Y/1000.0*yRange)).RGBA()

    r := float64(((r1+r2+r3+rc)/4) >> 8) / 255.
    g := float64(((g1+g2+g3+gc)/4) >> 8) / 255.
    b := float64(((b1+b2+b3+bc)/4) >> 8) / 255.

    return r,g,b
}

func triangulateImage(d sc.Delaunay, drawVertices, drawEdges, drawFaces bool) {
    file, err := os.Open("./deadpool.png")
    if err != nil {
        fmt.Println(err)
    }
    defer file.Close()

    img, err := png.Decode(file)
    if err != nil {
        fmt.Printf("%s: %v\n", "./kaktusfeige.png\n", err)
    }

    b := img.Bounds()
    //rgba := image.NewRGBA(b)
    //draw.Draw(rgba, b, img, b.Min, draw.Src)

    var imageSizeX float64 = 1000
    var imageSizeY float64 = 1000
    dc := gg.NewContext(int(imageSizeX), int(imageSizeY))

    // Background filling in white
    dc.SetRGB(1,1,1)
    dc.Clear()

    dc.SetRGB(1, 0, 0)

    for _, f := range d.Faces {

        v1 := d.Vertices[d.Edges[f.EEdge].VOrigin].Pos
        v2 := d.Vertices[d.Edges[d.Edges[f.EEdge].ENext].VOrigin].Pos
        v3 := d.Vertices[d.Edges[d.Edges[d.Edges[f.EEdge].ENext].ENext].VOrigin].Pos

        if drawFaces {
            dc.SetLineWidth(1.0)
            dc.SetRGB(calcFaceColor(img, b, v1, v2, v3))

            dc.MoveTo(v1.X, v1.Y)
            dc.LineTo(v2.X, v2.Y)
            dc.LineTo(v3.X, v3.Y)
            dc.LineTo(v1.X, v1.Y)
            dc.ClosePath()
            dc.Fill()

            dc.SetRGB(0.9902,0.9902,0.9902)

            dc.DrawLine(v1.X, v1.Y, v2.X, v2.Y)
            dc.Stroke()

            dc.DrawLine(v2.X, v2.Y, v3.X, v3.Y)
            dc.Stroke()

            dc.DrawLine(v3.X, v3.Y, v1.X, v1.Y)
            dc.Stroke()
            dc.Fill()

        }

        if drawEdges {
            dc.SetLineWidth(2.0)
            dc.SetRGB(calcEdgeColorFromVertices(img, b, v1, v2))
            dc.DrawLine(v1.X, v1.Y, v2.X, v2.Y)
            dc.Stroke()

            dc.SetRGB(calcEdgeColorFromVertices(img, b, v2, v3))
            dc.DrawLine(v2.X, v2.Y, v3.X, v3.Y)
            dc.Stroke()

            dc.SetRGB(calcEdgeColorFromVertices(img, b, v3, v1))
            dc.DrawLine(v3.X, v3.Y, v1.X, v1.Y)
            dc.Stroke()
            dc.Fill()
        }
    }


    if drawVertices {
        for _,v := range d.Vertices {

            if v == sc.EmptyV {
                continue
            }

            xRange := float64(b.Max.X-b.Min.X)
            yRange := float64(b.Max.Y-b.Min.Y)

            xRelPos := v.Pos.X/1000.0
            yRelPos := v.Pos.Y/1000.0

            r, g, b, _ := img.At(int(xRelPos*xRange), int(yRelPos*yRange)).RGBA()
            dc.SetRGB(float64(r>>8)/255., float64(g>>8)/255., float64(b>>8)/255.)

            dc.DrawCircle(v.Pos.X, v.Pos.Y, 2)
            dc.Fill()
        }
    }

    dc.SavePNG("triangulatedImage.png")

}

func drawImage(d sc.Delaunay) {
    var scale float64 = 1.0
    var imageSizeX float64 = 1000
    var imageSizeY float64 = 1000
    dc := gg.NewContext(int(imageSizeX), int(imageSizeY))

    // Background filling in white
    dc.SetRGB(1,1,1)
    dc.Clear()


    dc.SetLineWidth(1.0)

    for i := 1; i < 10; i++{

        x := float64(i)*100*scale
        y := float64(i)*100*scale

        dc.SetRGB(1, 0.5, 0.5)
        dc.DrawLine(0, y, imageSizeX, y)
        dc.Stroke()
        dc.DrawLine(x, 0, x, imageSizeY)
        dc.Stroke()

        dc.SetRGB(1, 0.0, 0.0)
        // X axis
        //dc.DrawString(strconv.Itoa(int(x)), x+10, imageSizeY-10)
        //// Y axis
        //dc.DrawString(strconv.Itoa(int(imageSizeY-y)), 10, y-10)

    }

    dc.SetLineWidth(2.0)
    for i,e := range d.Edges {

        if e == sc.EmptyE {
            continue
        }

        dc.SetRGB(0, 0, 0)

        v1 := d.Vertices[e.VOrigin].Pos
        v2 := d.Vertices[d.Edges[e.ETwin].VOrigin].Pos

        dc.DrawLine(v1.X*scale, v1.Y*scale, v2.X*scale, v2.Y*scale)
        dc.Stroke()

        //dc.SetRGB(0, 0, 1)
        //dc.DrawString(fmt.Sprintf("(%.1f, %.1f)", v1.X, v1.Y), v1.X, imageSizeY-v1.Y)
        //
        //dc.SetRGB(0, 0.5, 0)
        //middleP := v.Vector{(v1.X+v2.X)/2., (v1.Y+v2.Y)/2., 0}
        //ortho   := v.Vector{0,0,1}
        //crossP  := v.Cross(v.Sub(v1, v2), ortho)
        //crossP.Div(v.Length(crossP))
        //crossP.Mult(15.)
        //
        //middleP.Add(crossP)
        //
        i = i
        //s := fmt.Sprintf("(%d)", i)
        //dc.DrawStringAnchored(s, middleP.X, imageSizeY-middleP.Y, 0.5, 0.5)
    }

    dc.SetLineWidth(1.0)
    dc.SetRGB(1, 0, 0)
    for i,v := range d.Vertices {

        if v == sc.EmptyV {
            continue
        }

        dc.DrawCircle(v.Pos.X*scale, v.Pos.Y*scale, 2)
        dc.Fill()

        i = i
        //s := fmt.Sprintf("(%d)", i)
        //dc.DrawStringAnchored(s, v.Pos.X-10, imageSizeY-v.Pos.Y-10, 0.5, 0.5)
    }

    //dc.SetRGB(1, 1, 0)
    //dc.DrawCircle(432, imageSizeY-894, 5)
    //dc.DrawCircle(599, imageSizeY-532, 5)
    //dc.DrawCircle(501, imageSizeY-578, 5)
    //dc.Fill()

    dc.SavePNG("out_me.png")
}

func drawFgmImage(points []fgm.Point, triangulation *fgm.Triangulation) {
	
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

	dc.SavePNG("out.png")	
}

func triangulate(myPoints v.PointList, fgmPoints []fgm.Point, renderImage, profileMT, profileFgm bool) {
	
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
		if err != nil{
			fmt.Printf("Fogleman encountered an error: %v\n", err)
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
    
    if profileFgm || profileMT {
		profiling.Stop()
	}
	
	errFgm := triangulationFgm.Validate()
	errMT  := triangulationMT.Verify()
	if errFgm != nil {
		fmt.Printf("Fogleman encountered an error: %v\n", errFgm)
	}
	if errMT != nil {
		fmt.Printf("SweepCircle encountered an error: %v\n", errMT)
	}

    fmt.Printf("Triangulation (ms): %.8f, Fogleman (ms): %.8f\n", float64(binTime)/1000000.0, float64(binTimeFgm)/1000000.0)

	if renderImage {
		drawImage(triangulationMT)
		drawFgmImage(fgmPoints, triangulationFgm)
	}
	
}

func toFgmList(points v.PointList) []fgm.Point {
	newPoints := make([]fgm.Point, len(points), len(points))
	for i,v := range(points) {
		newPoints[i] = fgm.Point{v.X, v.Y}
	}
	return newPoints
}

func testUnknownProblemRandom(count int) (v.PointList, []fgm.Point) {
    fmt.Printf("===========================\n")
    fmt.Printf("=== test random\n")
    fmt.Printf("===========================\n")
    
    scale := 1000.0
    margin := 10.0
    
    var seed int64 = time.Now().UTC().UnixNano()
    seed = seed
    fmt.Fprintf(os.Stderr, "Seed: %v\n", seed)
    r := rand.New(rand.NewSource(1533982892382782961))
    var pointList v.PointList

    for i:= 0; i < count; i++ {
        pointList = append(pointList, v.Vector{r.Float64()*(scale-2*margin)+margin, r.Float64()*(scale-2*margin)+margin})
    }
    
	return pointList, toFgmList(pointList) 
}

func testCircle(count int) (v.PointList, []fgm.Point) {
    fmt.Printf("===========================\n")
    fmt.Printf("=== test Circle\n")
    fmt.Printf("===========================\n")

    var pointList v.PointList
    center := v.Vector{500, 500}
    radius := 350.

    for i:= 0; i < count; i++ {
        newI := v.DegToRad(float64(i)/float64(count)*float64(i))
        v := v.Vector{center.X + radius * math.Cos(float64(newI)), center.Y + radius * math.Sin(float64(newI))}
        
        pointList = append(pointList, v)
    }
    
    //radius = 450.
	//
    //for i:= 0; i < count; i++ {
    //    newI := v.DegToRad(float64(i)/float64(count)*360.)
    //    v := v.Vector{center.X + radius * math.Cos(float64(newI)), center.Y + radius * math.Sin(float64(newI))}
    //    pointList = append(pointList, v)
    //}

    //pointList = append(pointList, v.Vector{500, 500})
    
    return pointList, toFgmList(pointList) 
}

func testWaveCenterMirrored(count int) (v.PointList, []fgm.Point) {
    fmt.Printf("===========================\n")
    fmt.Printf("=== test wave center mirrored\n")
    fmt.Printf("===========================\n")
  
    var pointList v.PointList

    for i := 0; i <= count; i++ {

            newI := v.DegToRad(float64(i)/float64(count)*360.*5.0)

            v := v.Vector{float64(i)/float64(count)*900+50, math.Sin(newI)*450.+500}
            //fmt.Println(v)
            pointList = append(pointList, v)
    }

    return pointList, toFgmList(pointList) 
}

func testWave(count int) (v.PointList, []fgm.Point) {
    fmt.Printf("===========================\n")
    fmt.Printf("=== test wave\n")
    fmt.Printf("===========================\n")

    var pointList v.PointList

    for i := 0; i <= count; i++ {
            newI := v.DegToRad(float64(i)/float64(count)*360.)

            v := v.Vector{float64(i)/float64(count)*900+50, math.Sin(newI)*450.+500}
            //fmt.Println(v)
            pointList = append(pointList, v)
    }

    return pointList, toFgmList(pointList) 
}

func testTiltedGrid(count int, tiltAngle float64) (v.PointList, []fgm.Point) {
    fmt.Printf("===========================\n")
    fmt.Printf("=== test tilted grid\n")
    fmt.Printf("===========================\n")
    
    angle := v.DegToRad(tiltAngle)

    var pointList v.PointList

    for x := 0; x <= count; x++ {
        newX := float64(x)/float64(count) * 900. + 50.
        for y := 0; y <= count; y++ {
            newY := float64(y)/float64(count) * 900. + 50.
            tiltedX := (newX-500.) * math.Cos(angle) - (newY-500.) * math.Sin(angle)
            tiltedY := (newY-500.) * math.Cos(angle) + (newX-500.) * math.Sin(angle)
            tiltedX += 500.
            tiltedY += 500.
            v := v.Vector{tiltedX, tiltedY}
            pointList = append(pointList, v)
        }
    }

    return pointList, toFgmList(pointList) 
}

func main() {

	var myP v.PointList
	var fgmP []fgm.Point
    myP, fgmP = testUnknownProblemRandom(5) 
    triangulate(myP, fgmP, true, false, false)
    //d = testTiltedGrid(50, 0.0)
    //d = testTiltedGrid(50, 89.0)
    //d = testTiltedGrid(50, 45.0)
    //myP, fgmP = testCircle(100)
    //triangulate(myP, fgmP, false, false, false)
    //d = testWave()

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

}
