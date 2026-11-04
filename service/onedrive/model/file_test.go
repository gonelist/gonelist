package model

import (
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	db.Close()
	var err error
	db, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(1)
	initTable()
	code := m.Run()
	db.Close()
	os.Exit(code)
}

func testNode(t *testing.T) *FileNode {
	t.Helper()
	return &FileNode{ID: t.Name(), Name: "test.txt", Path: "/" + t.Name(), ParentID: "parent", LastModifyTime: time.Unix(1, 0)}
}

func TestInsert(t *testing.T) {
	node := testNode(t)
	if err := InsetFile(node); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for {
		got, err := Find(node.ID)
		if err == nil {
			if got.Path != node.Path {
				t.Fatal(got)
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatal(err)
		}
		time.Sleep(time.Millisecond)
	}
}

func TestFileLifecycle(t *testing.T) {
	node := testNode(t)
	if err := BatchInsertFile([]*FileNode{node}); err != nil {
		t.Fatal(err)
	}
	got, err := FindByPath(node.Path)
	if err != nil || got.ID != node.ID {
		t.Fatalf("find by path: %v, %v", got, err)
	}
	node.Name = "renamed.txt"
	if err := UpdateFile(node); err != nil {
		t.Fatal(err)
	}
	got, err = Find(node.ID)
	if err != nil || got.Name != node.Name {
		t.Fatalf("update: %v, %v", got, err)
	}
	nodes, err := FindByName(node.Name)
	if err != nil || len(nodes) != 1 {
		t.Fatalf("find by name: %v, %v", nodes, err)
	}
	nodes, err = GetChildrenByID(node.ParentID)
	if err != nil || len(nodes) == 0 {
		t.Fatalf("children: %v, %v", nodes, err)
	}
	if err := DeleteFile(node.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := Find(node.ID); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("deleted node: %v", err)
	}
}
