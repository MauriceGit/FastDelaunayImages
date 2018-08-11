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
        y := imageSizeY - float64(i)*100*scale

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

        dc.DrawLine(v1.X*scale, imageSizeY-v1.Y*scale, v2.X*scale, imageSizeY-v2.Y*scale)
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

        dc.DrawCircle(v.Pos.X*scale, imageSizeY-v.Pos.Y*scale, 2)
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

    dc.SavePNG("out.png")
}

func testUnknownProblem01() sc.Delaunay {
    fmt.Printf("===========================\n")
    fmt.Printf("=== test_unknown_problem_01\n")
    fmt.Printf("===========================\n")
    count := 10000
    var seed int64 = time.Now().UTC().UnixNano()
    seed = seed
    r := rand.New(rand.NewSource(4))
    var pointList v.PointList

    for i:= 0; i < count; i++ {
        v := v.Vector{r.Float64()*900+50, r.Float64()*900+50, 0}
        pointList = append(pointList, v)
    }


    return sc.Triangulate(pointList)
}

func testUnknownProblem02() sc.Delaunay {
    fmt.Printf("===========================\n")
    fmt.Printf("=== test_unknown_problem_02\n")
    fmt.Printf("===========================\n")
    count := 10
    r := rand.New(rand.NewSource(1525941446937387107))
    var pointList v.PointList

    for i:= 0; i < count; i++ {
        v := v.Vector{r.Float64()*900+50, r.Float64()*900+50, 0}
        pointList = append(pointList, v)
    }

    return sc.Triangulate(pointList)
}

func testUnknownProblem03() sc.Delaunay {
    fmt.Printf("===========================\n")
    fmt.Printf("=== test_unknown_problem_03\n")
    fmt.Printf("===========================\n")
    count := 2000
    r := rand.New(rand.NewSource(1525942373618049687))
    var pointList v.PointList

    for i:= 0; i < count; i++ {
        v := v.Vector{r.Float64()*900+50, r.Float64()*900+50, 0}
        pointList = append(pointList, v)
    }

    start := time.Now()
    delaunay := sc.Triangulate(pointList)
    binTime := time.Since(start).Nanoseconds()

    fmt.Printf("Triangulation in Milliseconds: %.8f\n", float64(binTime)/1000000000.0)

    return delaunay
}

func testUnknownProblemRandom(count int, scale, margin float64) sc.Delaunay {
    fmt.Printf("===========================\n")
    fmt.Printf("=== test_unknown_problem_random\n")
    fmt.Printf("===========================\n")
    //count := 5
    var seed int64 = time.Now().UTC().UnixNano()
    seed = seed
    fmt.Fprintf(os.Stderr, "Seed: %v\n", seed)
    r := rand.New(rand.NewSource(1533982892382782961))
    var pointList v.PointList

    for i:= 0; i < count; i++ {
        v := v.Vector{r.Float64()*(scale-2*margin)+margin, r.Float64()*(scale-2*margin)+margin, 0}
        pointList = append(pointList, v)
    }


    start := time.Now()
    delaunay := sc.Triangulate(pointList)
    binTime := time.Since(start).Nanoseconds()

    fmt.Printf("Triangulation in Milliseconds: %.8f\n", float64(binTime)/1000000.0)

    return delaunay
}

func testCircle() sc.Delaunay {
    fmt.Printf("===========================\n")
    fmt.Printf("=== test_unknown_problem_random\n")
    fmt.Printf("===========================\n")
    count := 1000
    var pointList v.PointList

    center := v.Vector{500, 500, 0}
    radius := 450.

    for i:= 0; i < count; i++ {
        newI := v.DegToRad(float64(i)/float64(count)*360.)
        v := v.Vector{center.X + radius * math.Cos(float64(newI)), center.Y + radius * math.Sin(float64(newI)), 0}
        pointList = append(pointList, v)
    }

    pointList = append(pointList, v.Vector{500, 500, 0})

    start := time.Now()
    delaunay := sc.Triangulate(pointList)
    binTime := time.Since(start).Nanoseconds()

    fmt.Printf("Triangulation in Milliseconds: %.8f\n", float64(binTime)/1000000.0)

    return delaunay
}



func testWaveCenterMirrored() sc.Delaunay {
    fmt.Printf("===========================\n")
    fmt.Printf("=== test_unknown_problem_random\n")
    fmt.Printf("===========================\n")
    count := 15
    var pointList v.PointList

    for i := 0; i <= count; i++ {

            newI := v.DegToRad(float64(i)/float64(count)*360.*5.0)

            v := v.Vector{float64(i)/float64(count)*900+50, math.Sin(newI)*450.+500, 0}
            //fmt.Println(v)
            pointList = append(pointList, v)
    }

    start := time.Now()
    delaunay := sc.Triangulate(pointList)
    binTime := time.Since(start).Nanoseconds()

    fmt.Printf("Triangulation in Milliseconds: %.8f\n", float64(binTime)/1000000.0)

    return delaunay
}

func testWave() sc.Delaunay {
    fmt.Printf("===========================\n")
    fmt.Printf("=== test_unknown_problem_random\n")
    fmt.Printf("===========================\n")
    count := 150
    var pointList v.PointList

    for i := 0; i <= count; i++ {

            newI := v.DegToRad(float64(i)/float64(count)*360.)

            v := v.Vector{float64(i)/float64(count)*900+50, math.Sin(newI)*450.+500, 0}
            //fmt.Println(v)
            pointList = append(pointList, v)
    }

    start := time.Now()
    delaunay := sc.Triangulate(pointList)
    binTime := time.Since(start).Nanoseconds()

    fmt.Printf("Triangulation in Milliseconds: %.8f\n", float64(binTime)/1000000.0)

    return delaunay
}

func testTiltedGrid(tiltAngle float64) sc.Delaunay {
    fmt.Printf("===========================\n")
    fmt.Printf("=== test_grid\n")
    fmt.Printf("===========================\n")
    count := 50
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
            v := v.Vector{tiltedX, tiltedY, 0}
            pointList = append(pointList, v)
        }
    }

    start := time.Now()
    delaunay := sc.Triangulate(pointList)
    binTime := time.Since(start).Nanoseconds()

    fmt.Printf("Triangulation in Milliseconds: %.8f\n", float64(binTime)/1000000.0)

    return delaunay
}

func main() {

    defer profile.Start(profile.CPUProfile).Stop()

    var d sc.Delaunay
    //d = testUnknownProblem03()
    //d = testUnknownProblemRandom(13000, 1000, 10)
    //d = testTiltedGrid(0.0)
    //d = testTiltedGrid(89.0)
    //d = testTiltedGrid(45.0)
    //d = testCircle()
    d = testWave()
    d = d

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




    //fmt.Println(d)
    fmt.Println(d.Verify())
    //fmt.Println(d)

    //triangulateImage(d, false, false, true)
    drawImage(d)
}
