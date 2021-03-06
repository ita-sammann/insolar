package artifactmanager

import (
	"os"
	"testing"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage/leveldb"
	"github.com/stretchr/testify/assert"
)

func prepareObjectDescriptorTest() (
	*leveldb.LevelLedger, *LedgerArtifactManager, *record.ObjectActivateRecord, *record.Reference,
) {
	if err := leveldb.DropDB(); err != nil {
		os.Exit(1)
	}
	ledger, _ := leveldb.InitDB()
	manager := LedgerArtifactManager{
		storer:   ledger,
		archPref: []record.ArchType{1},
	}
	rec := record.ObjectActivateRecord{Memory: record.Memory{1}}
	ref, _ := ledger.SetRecord(&rec)

	return ledger, &manager, &rec, ref
}

func TestObjectDescriptor_GetMemory(t *testing.T) {
	ledger, manager, objRec, objRef := prepareObjectDescriptorTest()
	amendRec := record.ObjectAmendRecord{NewMemory: record.Memory{2}}
	amendRef, _ := ledger.SetRecord(&amendRec)
	idx := index.ObjectLifeline{
		LatestStateRef: *amendRef,
	}
	ledger.SetObjectIndex(objRef, &idx)

	desc := ObjectDescriptor{
		manager:           manager,
		activateRecord:    objRec,
		latestAmendRecord: nil,
		lifelineIndex:     &idx,
	}
	mem, err := desc.GetMemory()
	assert.NoError(t, err)
	assert.Equal(t, record.Memory{1}, mem)

	desc = ObjectDescriptor{
		manager:           manager,
		activateRecord:    objRec,
		latestAmendRecord: &amendRec,
		lifelineIndex:     &idx,
	}
	mem, err = desc.GetMemory()
	assert.NoError(t, err)
	assert.Equal(t, record.Memory{2}, mem)
}

func TestObjectDescriptor_GetDelegates(t *testing.T) {
	ledger, manager, objRec, objRef := prepareObjectDescriptorTest()
	appendRec1 := record.ObjectAppendRecord{AppendMemory: record.Memory{2}}
	appendRec2 := record.ObjectAppendRecord{AppendMemory: record.Memory{3}}
	appendRef1, _ := ledger.SetRecord(&appendRec1)
	appendRef2, _ := ledger.SetRecord(&appendRec2)
	idx := index.ObjectLifeline{
		LatestStateRef: *objRef,
		AppendRefs:     []record.Reference{*appendRef1, *appendRef2},
	}
	ledger.SetObjectIndex(objRef, &idx)

	desc := ObjectDescriptor{
		manager:           manager,
		activateRecord:    objRec,
		latestAmendRecord: nil,
		lifelineIndex:     &idx,
	}

	appends, err := desc.GetDelegates()
	assert.NoError(t, err)
	assert.Equal(t, []record.Memory{{2}, {3}}, appends)
}
