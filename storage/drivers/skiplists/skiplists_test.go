package skiplists

import (
	"reflect"
	"testing"
)

func testSetup(n, lv int) {
	MaxNodes = 16
	MaxLevel = 4
}

func TestInsert(t *testing.T) {
	testcases := []struct {
		name        string
		maxNodes    int
		maxLevel    int
		insertNodes int
		wantErr     error
		wantLen     int
	}{
		{
			name:        "insert nodes < maxNode",
			maxNodes:    16,
			maxLevel:    4,
			insertNodes: 15,
			wantErr:     nil,
		},
		{
			name:        "insert nodes = maxNode",
			maxNodes:    16,
			maxLevel:    4,
			insertNodes: 16,
			wantErr:     nil,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			testSetup(tt.maxNodes, tt.maxLevel)
			list := New()
			for i := 1; i <= tt.insertNodes; i++ {
				if err := list.insert(i, "test-data"); err != nil {
					t.Fatalf("the error should be %v, but got %v", tt.wantErr, err)
				}
			}
			if list.len() != tt.insertNodes {
				t.Fatalf("the length should be %d, but got %d", tt.insertNodes, list.len())
			}
		})
	}
}

func TestInsertErr(t *testing.T) {
	testcases := []struct {
		name        string
		maxNodes    int
		maxLevel    int
		insertNodes int
		wantErr     error
		wantLen     int
	}{
		{
			name:        "insert nodes > maxNode",
			maxNodes:    16,
			maxLevel:    4,
			insertNodes: 17,
			wantErr:     ErrSkipListIsFull,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			testSetup(tt.maxNodes, tt.maxLevel)
			list := New()
			// insert nodes under the limitation
			for i := 1; i <= tt.maxNodes; i++ {
				if err := list.insert(i, "test-data"); err != nil {
					t.Fatalf("insert valid error %v", err)
				}
			}

			// insert additional nodes
			for i := tt.maxNodes + 1; i <= tt.insertNodes; i++ {
				if err := list.insert(i, "test-data"); err == nil {
					t.Fatalf("the error should be %v, but got %v", tt.wantErr, err)
				}
			}
		})
	}
}

func testList() *SkipList {
	maxNodes, maxLV := 16, 4
	testSetup(maxNodes, maxLV)
	return New()
}

func TestInsertWraperWithSequenceID(t *testing.T) {
	t.Log("Start testing the sequence id creation by insert wrapper...")
	list := testList()
	items := []string{"a", "b", "c", "d", "e"}
	prevID := 0
	for _, item := range items {
		id, err := list.Insert(item)
		if err != nil {
			t.Fatal("insert error", err)
		}

		// id should be greater than prevID and unique
		if id <= prevID {
			t.Fatal("id is less than and equal to prevID")
		}
		prevID = id
	}

	t.Log("show list")
	list.display()
	t.Log("... Passed")
}

func TestSearch(t *testing.T) {
	testcases := []struct {
		name     string
		in       []int
		key      int
		want     error
		wantData any
	}{
		{
			name:     "data was found",
			in:       []int{1, 3, 5, 6, 7, 8},
			key:      5,
			wantData: 5,
			want:     nil,
		},
		{
			name:     "data was not found",
			in:       []int{1, 3, 5, 6, 7, 8},
			key:      9,
			wantData: nil,
			want:     ErrSkipListDataNotFound,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			list := testList()
			for _, key := range tt.in {
				if err := list.insert(key, tt.wantData); err != nil {
					t.Fatal("insert error", err)
				}
			}

			data, err := list.search(tt.key)
			if err != tt.want {
				t.Fatalf("search error should be %v, but got %v", tt.want, err)
			}

			if !reflect.DeepEqual(data, tt.wantData) {
				t.Fatalf("search data should be %v, but got %v", tt.wantData, data)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	testcases := []struct {
		name string
		in   []int
		key  int
		want bool
	}{
		{
			name: "deleted",
			in:   []int{1, 3, 5, 6, 7, 8},
			key:  5,
			want: true,
		},
		{
			name: "deleted on non-exist data",
			in:   []int{1, 3, 5, 6, 7, 8},
			key:  9,
			want: false,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			list := testList()
			for _, key := range tt.in {
				if err := list.insert(key, "test-data"); err != nil {
					t.Fatal("insert error", err)
				}
			}

			ans := list.Delete(tt.key)
			if ans != tt.want {
				t.Fatalf("deleted should be %v, but got %v", tt.want, ans)
			}

			// using Search() to check the node is deleted
			if _, err := list.search(tt.key); err != ErrSkipListDataNotFound {
				t.Fatalf("search deleted value error should be %v, but got %v", ErrSkipListDataNotFound, err)
			}
		})
	}
}

func TestRangeQuery(t *testing.T) {
	testcases := []struct {
		name  string
		in    []int
		start int
		end   int
		want  []any
	}{
		{
			name:  "query - 0",
			in:    []int{1, 2, 3, 4, 5, 6},
			start: 1,
			end:   1,
			want:  []any{1},
		},
		{
			name:  "query - 1",
			in:    []int{1, 2, 3, 4, 5, 6},
			start: 1,
			end:   5,
			want:  []any{1, 2, 3, 4, 5},
		},
		{
			name:  "query - 2",
			in:    []int{1, 2, 3, 4, 5, 6},
			start: 2,
			end:   5,
			want:  []any{6},
		},
		{
			name:  "query - 3",
			in:    []int{1, 2, 3, 4, 5, 6},
			start: 3,
			end:   5,
			want:  []any{},
		},
		{
			name:  "query - 4",
			in:    []int{1, 2, 3, 4, 5, 6, 7},
			start: 1,
			end:   5,
			want:  []any{1, 2, 3, 4, 5},
		},
		{
			name:  "query - 5",
			in:    []int{1, 2, 3, 4, 5, 6, 7},
			start: 2,
			end:   5,
			want:  []any{6, 7},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			list := testList()
			for _, key := range tt.in {
				if err := list.insert(key, "test-data"); err != nil {
					t.Fatal("insert error", err)
				}
			}

			ans := list.Range(tt.start, tt.end)
			if !reflect.DeepEqual(ans, tt.want) {
				t.Fatalf("it should be %v, but got %v", tt.want, ans)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	testcases := []struct {
		name    string
		data    any
		newData any
	}{
		{
			name:    "test data updating",
			data:    "data - v1",
			newData: "data - v2",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			list := testList()
			id, err := list.Insert(tt.data)
			if err != nil {
				t.Fatal("insert error", err)
			}
			if err := list.Update(id, tt.newData); err != nil {
				t.Fatal("update error", err)
			}
			// check if data is updated
			data, err := list.search(id)
			if err != nil {
				t.Fatal("search error", err)
			}
			if !reflect.DeepEqual(data, tt.newData) {
				t.Fatalf("failed to update data, it should be %v, but got %v", tt.newData, data)
			}
		})
	}
}
