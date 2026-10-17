package gtfs

import (
	"context"
	"testing"
)

func TestInitializeGTFS(t *testing.T) {
	config := GTFSConfig{StaticUrl: "test_url", City: Madison}
	testGTFS, err := NewGTFS(config)

	if err != nil {
		t.Fatalf("failed to initialize gtfs %v", err)
	}

	if testGTFS.Config.City != Madison {
		t.Fatalf("unexpected city %s", testGTFS.Config.City)
	}

	if testGTFS.Config.StaticUrl != "test_url" {
		t.Fatalf("unexpected static url %s", testGTFS.Config.StaticUrl)
	}
}

func TestGTFSGetStaticFIleName(t *testing.T) {
	config := GTFSConfig{StaticUrl: "test_url", City: Madison}
	testGTFS, err := NewGTFS(config)

	if err != nil {
		t.Fatalf("failed to initialize gtfs %v", err)
	}

	name := testGTFS.GetStaticFileName()

	if name != "madison-gtfs.zip" {
		t.Fatalf("unexpected file name %s, expected madison-gtfs.zip", name)
	}
}

func TestGTFSGetLocalDataDir(t *testing.T) {
	config := GTFSConfig{StaticUrl: "test_url", City: Madison}
	testGTFS, err := NewGTFS(config)

	if err != nil {
		t.Fatalf("failed to initialize gtfs %v", err)
	}

	dir, err := testGTFS.GetLocalDataDir()

	if err != nil {
		t.Fatalf("failed to get local data directory %v", err)
	}

	if dir == "" {
		t.Fatalf("failed to get local data directory, returned null %v", err)
	}
}

func TestGTFSDownloadStatic(t *testing.T) {
	config := GTFSConfig{StaticUrl: "https://transitdata.cityofmadison.com/GTFS/mmt_gtfs.zip", City: Madison}
	testGTFS, err := NewGTFS(config)

	if err != nil {
		t.Fatalf("failed to initialize gtfs %v", err)
	}

	ctx := context.Background()
	err = testGTFS.DownloadStatic(ctx)

	if err != nil {
		t.Fatalf("failed to download static file %v", err)
	}

	if testGTFS.StaticHash == "" {
		t.Fatalf("failed to hash file %v", err)
	}
}

func TestGTFSParseStatic(t *testing.T) {
	config := GTFSConfig{StaticUrl: "https://transitdata.cityofmadison.com/GTFS/mmt_gtfs.zip", City: Madison}
	testGTFS, err := NewGTFS(config)

	if err != nil {
		t.Fatalf("failed to initialize gtfs %v", err)
	}

	err = testGTFS.ParseStatic()

	if err != nil {
		t.Fatalf("failed to parse static file %v", err)
	}

	if testGTFS.Agency == nil {
		t.Fatalf("expected agency not to be null %v", err)
	}

	if testGTFS.Agency.Id == "" || testGTFS.Agency.Name == "" {
		t.Fatalf("expected id and name not to be null %v", err)
	}

	if len(testGTFS.Routes) == 0 {
		t.Fatalf("expected Routes to contain routes %v", err)
	}

	for _, route := range testGTFS.Routes {
		if route.Id == "" || route.AgencyId == "" {
			t.Fatalf("Route must contain id and agency id %v", err)
		}
	}

	if testGTFS.FeedInfo == nil {
		t.Fatalf("expected FeedInfo not to be null %v", err)
	}

	if len(testGTFS.Stops) == 0 {
		t.Fatalf("expected Stops to contain stops %v", err)
	}

	if len(testGTFS.Shapes) == 0 {
		t.Fatalf("expected Stops to contain stops %v", err)
	}

	if len(testGTFS.Calendars) == 0 {
		t.Fatalf("expected Stops to contain stops %v", err)
	}

	if len(testGTFS.Trips) == 0 {
		t.Fatalf("expected Stops to contain stops %v", err)
	}

}
