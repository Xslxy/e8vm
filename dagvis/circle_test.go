package dagvis

import (
	"testing"

	"fmt"
	"reflect"
	"sort"
)

func TestMinCircle(t *testing.T) {
	o := func(nodes map[string][]string, circExpect []string) {
		m, err := initMap(&Graph{Nodes: nodes})
		if err != nil {
			t.Errorf("init map %v, got error: %v", nodes, err)
			return
		}
		res := minCircle(m.Nodes)

		var circGot []string
		for _, node := range res {
			circGot = append(circGot, node.Name)
		}
		sort.Strings(circGot)

		if !reflect.DeepEqual(circGot, circExpect) {
			t.Errorf("min circle of %v, got %v, expect %v",
				nodes, circGot, circExpect,
			)
		}
	}

	o(map[string][]string{}, nil)
	o(map[string][]string{"a": {}}, nil)

	o(map[string][]string{
		"1": {"2"},
		"2": {"1"},
	}, []string{"1", "2"})

	o(map[string][]string{
		"1": {"2"},
		"2": {"3"},
		"3": {"1"},
	}, []string{"1", "2", "3"})

	o(map[string][]string{
		"1": {"2"},
		"2": {"3"},
		"3": {},
	}, nil)

	o(map[string][]string{
		"1": {"2"},
		"2": {"3"},
		"3": {"1", "4"},
		"4": {"5"},
		"5": {"1"},
	}, []string{"1", "2", "3"})

	o(map[string][]string{
		"1": {"2"},
		"2": {"3", "5"},
		"3": {"4"},
		"4": {"1"},
		"5": {"1"},
	}, []string{"1", "2", "5"})

	o(map[string][]string{
		"1":  {"2"},
		"2":  {"3", "6"},
		"3":  {"4"},
		"4":  {"5"},
		"5":  {"1"},
		"6":  {"1"},
		"7":  {"1", "2", "3", "4", "5", "6"},
		"8":  {"1", "2", "3", "4", "5", "6"},
		"9":  {"1", "2", "3", "4", "5", "6"},
		"10": {"1", "2", "3", "4", "5", "6"},
		"11": {"1", "2", "3", "4", "5", "6"},
	}, []string{"1", "2", "6"})

	nodes := map[string][]string{
		"1": {"3", "4"},
		"3": {"6", "4"},
		"4": {"3", "8"},
		"6": {},
		"8": {"1"},
	}
	// add a bunch of empty nodes
	for i := 0; i < 100; i++ {
		k := fmt.Sprintf(":%d", i)
		nodes[k] = []string{}
	}

	o(nodes, []string{"3", "4"})

	nodes = map[string][]string{
		"1": {"2"},
		"2": {"3"},
		"3": {"1"},
	}

	for i := 0; i < 5; i++ {
		k1 := fmt.Sprintf("a%d", i)
		k2 := fmt.Sprintf("b%d", i)
		nodes["1"] = append(nodes["1"], k1)
		nodes[k1] = []string{k2}
		nodes[k2] = []string{"2"}
	}

	for i := 0; i < 5; i++ {
		k1 := fmt.Sprintf("c%d", i)
		k2 := fmt.Sprintf("d%d", i)
		nodes["2"] = append(nodes["2"], k1)
		nodes[k1] = []string{k2}
		nodes[k2] = []string{"3"}
	}

	for i := 0; i < 5; i++ {
		k1 := fmt.Sprintf("e%d", i)
		k2 := fmt.Sprintf("f%d", i)
		nodes["3"] = append(nodes["3"], k1)
		nodes[k1] = []string{k2}
		nodes[k2] = []string{"1"}
	}

	// because iterations on maps are random, we do it for several times so
	// that it will iterate with different random permutations.
	for i := 0; i < 5; i++ {
		// TODO(YsGH): uncomment this when the implementation is fixed.
		// o(nodes, []string{"1", "2", "3"})
	}
}
