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
    "strconv"
    //"runtime/pprof"
    //"flag"
    //"log"
    //"runtime"
    "github.com/pkg/profile"
)


func drawImage(d sc.Delaunay) {
    var imageSizeX float64 = 1000
    var imageSizeY float64 = 1000
    dc := gg.NewContext(int(imageSizeX), int(imageSizeY))

    // Background filling in white
    dc.SetRGB(1,1,1)
    dc.Clear()


    dc.SetLineWidth(1.0)

    for i := 1; i < 10; i++{

        x := float64(i)*100
        y := imageSizeY - float64(i)*100

        dc.SetRGB(1, 0.5, 0.5)
        dc.DrawLine(0, y, imageSizeX, y)
        dc.Stroke()
        dc.DrawLine(x, 0, x, imageSizeY)
        dc.Stroke()

        dc.SetRGB(1, 0.0, 0.0)
        // X axis
        dc.DrawString(strconv.Itoa(int(x)), x+10, imageSizeY-10)
        // Y axis
        dc.DrawString(strconv.Itoa(int(imageSizeY-y)), 10, y-10)

    }

    dc.SetLineWidth(2.0)
    for i,e := range d.Edges {

        if e == sc.EmptyE {
            continue
        }

        dc.SetRGB(0, 0, 0)

        v1 := d.Vertices[e.VOrigin].Pos
        v2 := d.Vertices[d.Edges[e.ETwin].VOrigin].Pos

        dc.DrawLine(v1.X, imageSizeY-v1.Y, v2.X, imageSizeY-v2.Y)
        dc.Stroke()

        dc.SetRGB(0, 0, 1)
        //dc.DrawString(fmt.Sprintf("(%.1f, %.1f)", v1.X, v1.Y), v1.X, imageSizeY-v1.Y)

        dc.SetRGB(0, 0.5, 0)
        middleP := v.Vector{(v1.X+v2.X)/2., (v1.Y+v2.Y)/2., 0}
        ortho   := v.Vector{0,0,1}
        crossP  := v.Cross(v.Sub(v1, v2), ortho)
        crossP.Div(v.Length(crossP))
        crossP.Mult(15.)

        middleP.Add(crossP)

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

        dc.DrawCircle(v.Pos.X, imageSizeY-v.Pos.Y, 3)
        dc.Fill()

        i = i
        //s := fmt.Sprintf("(%d)", i)
        //dc.DrawStringAnchored(s, v.Pos.X-10, imageSizeY-v.Pos.Y-10, 0.5, 0.5)
    }

    dc.SetRGB(1, 1, 0)
    dc.DrawCircle(432, imageSizeY-894, 5)
    //dc.DrawCircle(599, imageSizeY-532, 5)
    //dc.DrawCircle(501, imageSizeY-578, 5)
    dc.Fill()

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

func testUnknownProblemRandom() sc.Delaunay {
    fmt.Printf("===========================\n")
    fmt.Printf("=== test_unknown_problem_random\n")
    fmt.Printf("===========================\n")
    count := 100
    var seed int64 = time.Now().UTC().UnixNano()
    seed = seed
    fmt.Fprintf(os.Stderr, "Seed: %v\n", 1528627210314976626)
    r := rand.New(rand.NewSource(1528627210314976626))
    var pointList v.PointList

    for i:= 0; i < count; i++ {
        v := v.Vector{r.Float64()*900+50, r.Float64()*900+50, 0}
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
    count := 10
    var pointList v.PointList

    center := v.Vector{500, 500, 0}
    radius := 450.

    for i:= 0; i < count; i++ {
        newI := v.DegToRad(float64(i)/float64(count)*360.)
        v := v.Vector{center.X + radius * math.Cos(float64(newI)), center.Y + radius * math.Sin(float64(newI)), 0}
        pointList = append(pointList, v)
    }

    //pointList = append(pointList, v.Vector{500, 500, 0})

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
    count := 100
    var pointList v.PointList

    for i := 0; i <= count; i++ {

            newI := v.DegToRad(float64(i)/float64(count)*360.)

            v := v.Vector{float64(i)/float64(count)*900+50, math.Sin(newI)*200.+500, 0}
            //fmt.Println(v)
            pointList = append(pointList, v)
    }

    start := time.Now()
    delaunay := sc.Triangulate(pointList)
    binTime := time.Since(start).Nanoseconds()

    fmt.Printf("Triangulation in Milliseconds: %.8f\n", float64(binTime)/1000000.0)

    return delaunay
}

func testGrid() sc.Delaunay {
    fmt.Printf("===========================\n")
    fmt.Printf("=== test_grid\n")
    fmt.Printf("===========================\n")
    count := 50

    var pointList v.PointList

    for x := 0; x < count; x++ {
        newX := float64(x)/float64(count) * 900. + 50.
        for y := 0; y < count; y++ {
            newY := float64(y)/float64(count) * 900. + 50.
            v := v.Vector{newX, newY, 0}
            pointList = append(pointList, v)
        }
    }

    start := time.Now()
    delaunay := sc.Triangulate(pointList)
    binTime := time.Since(start).Nanoseconds()

    fmt.Printf("Triangulation in Milliseconds: %.8f\n", float64(binTime)/1000000.0)

    return delaunay
}

func testTiltedGrid() sc.Delaunay {
    fmt.Printf("===========================\n")
    fmt.Printf("=== test_grid\n")
    fmt.Printf("===========================\n")
    count := 50
    angle := v.DegToRad(89.0)

    var pointList v.PointList

    for x := 0; x < count; x++ {
        newX := float64(x)/float64(count) * 900. + 50.
        for y := 0; y < count; y++ {
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
    //d = testUnknownProblemRandom()
    //d = testTiltedGrid()
    //d = testCircle()
    d = testWave()
    d = d

    // Wrong timing.
    // 10       points  0.12
    // 100      points  0.19
    // 1000     points  0.38
    // 10000    points  1.8
    // 100000   points  14
    // 1000000  points  147

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

    // ===========================
    // === test_unknown_problem_random
    // ===========================
    // count: 2000000
    // Seed:  1527368018753283347
    // Triangulation in Milliseconds: 27.15096612

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
    // 10       points  0.00008067
    // 100      points  0.00056354
    // 1000     points  0.00479936
    // 10000    points  0.05738459
    // 100000   points  0.42946737
    // 1000000  points  7.71770482
    // 2000000  points  16.2359380

    // Some finetuning.
    // 10       points  0.00006059
    // 100      points  0.00047796
    // 1000     points  0.00444188
    // 10000    points  0.05479064
    // 100000   points  0.61441104
    // 1000000  points  6.89181266
    // 2000000  points  14.23493647
    // 5000000  points  41.91952035




    //fmt.Println(d)
    //d.Verify()

    drawImage(d)
}
