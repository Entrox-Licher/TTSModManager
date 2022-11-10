package tests

import (
	"ModCreator/file"
	"ModCreator/mod"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	expectedStr       = []string{"SaveName", "Date", "VersionNumber", "GameMode", "GameType", "GameComplexity", "Table", "Sky", "Note", "LuaScript", "LuaScriptState", "XmlUI"}
	expectedObj       = []string{"TabStates", "MusicPlayer", "Grid", "Lighting", "Hands", "ComponentTags", "Turns"}
	expectedObjArr    = []string{"CameraStates", "DecalPallet", "CustomUIAssets", "SnapPoints", "Decals"}
	expectedObjStates = "ObjectStates"
)

func TestAllReverseThenBuild(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join("testdata", "e2e", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		_, filename := filepath.Split(path)
		testname := filename[:len(filename)-len(filepath.Ext(path))]
		denyList := []string{
			"small_lua", // currently tries to bundle it
			"long_lua",  // currently tries to bundle it
		}

		t.Run(testname, func(t *testing.T) {
			for _, f := range denyList {
				if f == testname {
					return
				}
			}
			j, err := file.ReadRawFile(path)
			if err != nil {
				t.Fatalf("Error parsing %s : %v", path, err)
			}
			want, err := file.ReadRawFile(path)
			if err != nil {
				t.Fatalf("Error parsing %s : %v", path, err)
			}
			modsettings := NewFF()
			finalOutput := NewFF()
			objsAndLua := NewFF()

			r := mod.Reverser{
				ModSettingsWriter: modsettings,
				LuaWriter:         objsAndLua,
				ObjWriter:         objsAndLua,
				ObjDirCreeator:    objsAndLua,
				RootWrite:         finalOutput,
			}
			err = r.Write(j)
			if err != nil {
				t.Fatalf("Error reversing : %v", err)
			}

			objsAndLua.DebugFileNames(t.Logf)
			finalOutput.DebugFileNames(t.Logf)
			reversedConfig, _ := finalOutput.ReadObj("config.json")
			t.Logf("%v\n", reversedConfig)

			m := &mod.Mod{
				Lua:         objsAndLua,
				Modsettings: modsettings,
				Objs:        objsAndLua,
				Objdirs:     objsAndLua,
				RootRead:    finalOutput,
				RootWrite:   finalOutput,
			}
			err = m.GenerateFromConfig()
			if err != nil {
				t.Fatalf("generateMod(<config>) : %v\n", err)
			}
			err = m.Print()
			if err != nil {
				t.Fatalf("printMod(...) : %v", err)
			}
			got, err := finalOutput.ReadObj("output.json")
			if err != nil {
				t.Fatalf("output.json not parsed : %v", err)
			}

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
		})

	}
}