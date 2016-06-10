package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	//	"strconv"
)

type anything interface{}

type Edge struct {
	Eid         string
	Label       string
	Tail        string
	Head        string
	Description string
	Props       map[string]anything
}

type Vertex struct {
	Vid   string
	Vname string
	Label string
	In    []string
	Out   []string
	Props map[string]anything
}

type VertexIndex map[string]*Vertex
type EdgeIndex map[string]*Edge

type Graph struct {
	Edges    []*Edge
	Vertices []*Vertex
	Vi       VertexIndex
	Ei       EdgeIndex
	AutoVId  int
	AutoEId  int
}

type Query struct {
	G     *Graph
	Verts []*Vertex
	Root  string
}

// Returns empty graph
func CreateGraph() Graph {
	var g Graph
	g.AutoVId = 1
	g.AutoEId = 1
	g.Vi = make(map[string]*Vertex)
	g.Ei = make(map[string]*Edge)
	return g
}

// Returns a new vertex
func NewVertex() Vertex {
	var v Vertex
	return v
}

// Setter
func (v *Vertex) SetVertexValues(vname string, label string, in []string, out []string, props map[string]anything) {
	v.Vname = vname
	v.Label = label
	v.In = in
	v.Out = out
	v.Props = props
}

// Returns a new edge
func NewEdge() Edge {
	var e Edge
	return e
}

// Setter
func (e *Edge) SetEdgeValues(label string, head string, tail string, props map[string]anything) {
	e.Label = label
	e.Head = head
	e.Tail = tail
	e.Props = props
	e.Description = fmt.Sprintf("%s------- %s ---->>>-----%s ", tail, label, head)
	//fmt.Println(e.Description)
}

// appends vertices to the graph.Vertices
func (g *Graph) AddVertices(v []Vertex) error {
	return errors.New("added vertices successfully")
}

// appends vertex to the graph.Vertices
func (g *Graph) AddVertex(v *Vertex) error {
	// v.Vid = strconv.Itoa(g.AutoVId)
	v.Vid = fmt.Sprintf("V%d", g.AutoVId)
	g.Vertices = append(g.Vertices, v)
	g.Vi[v.Vname] = v
	g.AutoVId++

	return errors.New("added vertex successfully")
}

// appends edges to the graph.Edges
func (g *Graph) AddEdges(e []Edge) error {
	return errors.New("error adding edges")
}

// appends edge to the graph.Edges
func (g *Graph) AddEdge(e Edge) error {
	// e.Eid = strconv.Itoa(g.AutoEId)
	e.Eid = fmt.Sprintf("E%d", g.AutoEId)
	g.AutoEId++
	g.Edges = append(g.Edges, &e)
	g.Ei[e.Eid] = &e

	if vh, ok := g.findVertexByName(e.Head); ok { // check if vertex is present
		// fmt.Println("found VertexHead")
		vh.In = append(vh.In, e.Eid)
	}
	if vt, ok := g.findVertexByName(e.Tail); ok {
		// fmt.Println("found VertexTail")
		vt.Out = append(vt.Out, e.Eid)
	}

	return errors.New("cannot add edge")
}

// lookup table of id -> *Vertex
// returns if found --> *Vertex and true
//		   not found --> nil and false
func (g *Graph) findVertexByName(name string) (*Vertex, bool) {

	v, ok := g.Vi[name]
	if ok {
		// fmt.Println("vertex exists")
		return v, true
	}
	return nil, false
}

// returns pointer to the edge
func (g *Graph) findEdgeById(eid string) (*Edge, bool) {
	if e, ok := g.Ei[eid]; ok {
		//fmt.Println("edge exists")
		return e, true
	}
	return nil, false
}

// takes name of the node ; returns Query structure for cascading
func (g *Graph) Query(s string) *Query {
	var q Query
	var arrVerts = make([]*Vertex, 0)
	v, ok := g.findVertexByName(s)
	if ok {
		fmt.Println("found ", v.Vname)
	}
	arrVerts = append(arrVerts, v)
	q.Verts = arrVerts
	q.G = g
	q.Root = s
	return &q
}

func (q *Query) out(l string) *Query {

	var result = make([]*Vertex, 0)
	var output = make([]*Vertex, 0)

	result = q.Verts
	for _, val := range result {
		for _, oe := range val.Out {
			if ed, ok := q.G.findEdgeById(oe); ok {
				if ed.Label == l {
					fmt.Println(val.Vname, "  ", ed.Label, "   ", ed.Head)
					nv, _ := q.G.findVertexByName(ed.Head)
					output = append(output, nv)
				}
			}
		}
	}
	q.Verts = output
	return q
}

func (q *Query) value() {
	for _, v := range q.Verts {
		fmt.Printf("%+v\n", *v)
	}
}

func (q *Query) in(l string) *Query {
	var result = make([]*Vertex, 0)
	var input = make([]*Vertex, 0)

	result = q.Verts
	for _, val := range result {
		for _, ie := range val.In {
			if ed, ok := q.G.findEdgeById(ie); ok {
				// fmt.Println("found edge")
				if ed.Label == l {
					fmt.Println(ed.Tail, "  ", ed.Label, "   ", val.Vname)
					nv, _ := q.G.findVertexByName(ed.Tail)
					input = append(input, nv)
				}
			}
		}
	}
	q.Verts = input
	return q
}

func (q *Query) all() {

	arrVerts := q.Verts

	for _, ver := range arrVerts {
		ineds := ver.In
		outeds := ver.Out
		fmt.Println("all in <-- relations of ", ver.Vname)
		for _, in := range ineds {
			if ed, err := q.G.findEdgeById(in); err {
				fmt.Println(ed.Label)
			}

		}
		fmt.Println("all out --> relations of ", ver.Vname)
		for _, out := range outeds {
			if ed, err := q.G.findEdgeById(out); err {
				fmt.Println(ed.Label)
			}
		}

	}

	// ver, ok := q.G.findVertexByName(q.Root)

	// if ok {
	// 	ineds := ver.In
	// 	outeds := ver.Out
	// 	fmt.Println("all relations of ", q.Root)
	// 	for _, in := range ineds {
	// 		if ed, err := q.G.findEdgeById(in); err {
	// 			fmt.Println(ed.Label)
	// 		}

	// 	}
	// 	for _, out := range outeds {
	// 		if ed, err := q.G.findEdgeById(out); err {
	// 			fmt.Println(ed.Label)
	// 		}
	// 	}
	// } else {
	// 	fmt.Println("invalid input")
	// }
}

// Note : when declaring composite types in structure you need to mention the type
// for eg: below mentioning []string before array, map[string]anything
func main() {
	g := CreateGraph()
	v1 := NewVertex()
	v2 := NewVertex()
	v3 := NewVertex()
	v4 := NewVertex()
	v5 := NewVertex()
	v6 := NewVertex()
	v7 := NewVertex()
	v8 := NewVertex()
	v9 := NewVertex()
	v10 := NewVertex()
	v11 := NewVertex()
	v12 := NewVertex()

	e1 := NewEdge()
	e2 := NewEdge()
	e3 := NewEdge()
	e4 := NewEdge()
	e5 := NewEdge()
	e6 := NewEdge()
	e7 := NewEdge()
	e8 := NewEdge()
	e9 := NewEdge()
	e10 := NewEdge()
	e11 := NewEdge()
	e12 := NewEdge()
	e13 := NewEdge()
	e14 := NewEdge()
	e15 := NewEdge()
	e16 := NewEdge()
	e17 := NewEdge()

	// v1.SetVertexValues("Geekskool", []string{"Programming", "organisation"}, nil, nil, map[string]anything{"address": "indiranagar"})
	// v2.SetVertexValues("Santosh Rajan", []string{"Mentor", "Programmer"}, nil, nil, map[string]anything{"stars": 700, "project": "lispyScript"})
	// v3.SetVertexValues("Jayesh", []string{"Go Developer"}, nil, nil, map[string]anything{"address": "btm Layout", "project": "graph database"})

	v1.SetVertexValues("saturn", "titan", nil, nil, map[string]anything{"age": 10000})
	v2.SetVertexValues("sky", "location", nil, nil, map[string]anything{})
	v3.SetVertexValues("sea", "location", nil, nil, map[string]anything{})
	v4.SetVertexValues("jupiter", "god", nil, nil, map[string]anything{"age": 5000})
	v5.SetVertexValues("neptune", "god", nil, nil, map[string]anything{"age": 4500})
	v6.SetVertexValues("hercules", "demigod", nil, nil, map[string]anything{"age": 30})
	v7.SetVertexValues("alcmene", "human", nil, nil, map[string]anything{"age": 45})
	v8.SetVertexValues("pluto", "god", nil, nil, map[string]anything{"age": 4000})
	v9.SetVertexValues("nemean", "monster", nil, nil, map[string]anything{})
	v10.SetVertexValues("hydra", "monster", nil, nil, map[string]anything{})
	v11.SetVertexValues("cerberus", "monster", nil, nil, map[string]anything{})
	v12.SetVertexValues("tartarus", "location", nil, nil, map[string]anything{})

	g.AddVertex(&v1)
	g.AddVertex(&v2)
	g.AddVertex(&v3)
	g.AddVertex(&v4)
	g.AddVertex(&v5)
	g.AddVertex(&v6)
	g.AddVertex(&v7)
	g.AddVertex(&v8)
	g.AddVertex(&v9)
	g.AddVertex(&v10)
	g.AddVertex(&v11)
	g.AddVertex(&v12)

	// e1.SetEdgeValues([]string{"FounderOf", "CEOof"}, v1.Vid, v2.Vid, nil)
	// e2.SetEdgeValues([]string{"studenOf", "internOf"}, v2.Vid, v3.Vid, nil)
	//e4.SetEdgeValues(label, head, tail, props)
	e1.SetEdgeValues("father", v1.Vname, v4.Vname, nil)
	e2.SetEdgeValues("lives", v2.Vname, v4.Vname, map[string]anything{"reason": "loves fresh breezes"})
	e3.SetEdgeValues("father", v4.Vname, v6.Vname, nil)
	e4.SetEdgeValues("brother", v4.Vname, v5.Vname, nil)
	e5.SetEdgeValues("brother", v5.Vname, v4.Vname, nil)
	e6.SetEdgeValues("lives", v3.Vname, v5.Vname, map[string]anything{"reason": "loves waves"})
	e7.SetEdgeValues("mother", v7.Vname, v6.Vname, nil)
	e8.SetEdgeValues("brother", v8.Vname, v5.Vname, nil)
	e9.SetEdgeValues("brother", v5.Vname, v8.Vname, nil)
	e10.SetEdgeValues("brother", v4.Vname, v8.Vname, nil)
	e11.SetEdgeValues("brother", v8.Vname, v4.Vname, nil)
	e12.SetEdgeValues("battled", v9.Vname, v6.Vname, map[string]anything{"time": "1", "place": []float64{38.1, 23.7}})
	e13.SetEdgeValues("battled", v10.Vname, v6.Vname, map[string]anything{"time": "2", "place": []float64{37.7, 24}})
	e14.SetEdgeValues("battled", v11.Vname, v6.Vname, map[string]anything{"time": "12", "place": []float64{39, 22}})
	e15.SetEdgeValues("pet", v11.Vname, v8.Vname, nil)
	e16.SetEdgeValues("lives", v12.Vname, v11.Vname, nil)
	e17.SetEdgeValues("lives", v12.Vname, v8.Vname, map[string]anything{"reason": "fear of death"})

	g.AddEdge(e1)
	g.AddEdge(e2)
	g.AddEdge(e3)
	g.AddEdge(e4)
	g.AddEdge(e5)
	g.AddEdge(e6)
	g.AddEdge(e7)
	g.AddEdge(e8)
	g.AddEdge(e9)
	g.AddEdge(e10)
	g.AddEdge(e11)
	g.AddEdge(e12)
	g.AddEdge(e13)
	g.AddEdge(e14)
	g.AddEdge(e15)
	g.AddEdge(e16)
	g.AddEdge(e17)

	//data1, _ := json.MarshalIndent(g, "", "    ")
	f, ok := os.OpenFile("out.json", os.O_RDWR, 0666)

	if ok != nil {
		fmt.Println("file opening error")
	}

	json.NewEncoder(f).Encode(g)
	// g.outs()
	// g.ines()
	// v, k := g.findVertexByName("hercules")
	// v, k := g.findVertexByName("pluto")
	// v, k := g.findVertexByName("jupiter")
	// v, k := g.findVertexByName("neptune")

	// if k {
	// 	fmt.Println("found ", v.Vname)
	// }
	//	v.all()
	//v.out(&g, "father").out(&g, "father")
	//v.out(&g, "father").out(&g, "lives")
	// v.out(&g, "father").out(&g, "brother").out(&g, "lives")
	// v.out(&g, "father").out(&g, "brother").out(&g, "brother").out(&g, "pet").out(&g, "lives")
	// v.out(&g, "pet").out(&g, "lives")
	// v.out(&g, "battled")
	// v.ines(&g, "brother").out(&g, "brother").ines(&g, "father").out(&g, "mother")
	// v.out(&g, "brother")
	// g.Query("jupiter").in("father").out("battled").in("pet").out("lives")
	// g.Query("hydra").in("battled")
	// g.Query("tartarus").in("lives").in("pet").out("lives")
	// g.Query("hercules").out("battled").out("lives").all()
	g.Query("jupiter").out("lives").all()
	//	g.Query("jupiter").out("brother").out("lives").value()
}
