package transaction

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// const BRC74Hex = "fe8a6a0c000c04fde80b0011774f01d26412f0d16ea3f0447be0b5ebec67b0782e321a7a01cbdf7f734e30fde90b02004e53753e3fe4667073063a17987292cfdea278824e9888e52180581d7188d8fdea0b025e441996fc53f0191d649e68a200e752fb5f39e0d5617083408fa179ddc5c998fdeb0b0102fdf405000671394f72237d08a4277f4435e5b6edf7adc272f25effef27cdfe805ce71a81fdf50500262bccabec6c4af3ed00cc7a7414edea9c5efa92fb8623dd6160a001450a528201fdfb020101fd7c010093b3efca9b77ddec914f8effac691ecb54e2c81d0ab81cbc4c4b93befe418e8501bf01015e005881826eb6973c54003a02118fe270f03d46d02681c8bc71cd44c613e86302f8012e00e07a2bb8bb75e5accff266022e1e5e6e7b4d6d943a04faadcf2ab4a22f796ff30116008120cafa17309c0bb0e0ffce835286b3a2dcae48e4497ae2d2b7ced4f051507d010a00502e59ac92f46543c23006bff855d96f5e648043f0fb87a7a5949e6a9bebae430104001ccd9f8f64f4d0489b30cc815351cf425e0e78ad79a589350e4341ac165dbe45010301010000af8764ce7e1cc132ab5ed2229a005c87201c9a5ee15c0f91dd53eff31ab30cd4"
// const BRC74Root = "57aab6e6fb1b697174ffb64e062c4728f2ffd33ddcfa02a43b64d8cd29b483b4"
// const BRC74TXID1 = "304e737fdfcb017a1a322e78b067ecebb5e07b44f0a36ed1f01264d2014f7711"
// const BRC74TXID2 = "d888711d588021e588984e8278a2decf927298173a06737066e43f3e75534e00"
// const BRC74TXID3 = "98c9c5dd79a18f40837061d5e0395ffb52e700a2689e641d19f053fc9619445e"

// func TestMerklePath_ParseHex(t *testing.T) {
// 	t.Parallel()

// 	t.Run("parses from hex", func(t *testing.T) {
// 		path := MerklePath{}
// 		err := path.ParseHex(BRC74Hex)
// 		assert.NoError(t, err)
// 		assert.Equal(t, BRC74JSON.Path, path.Path)
// 	})
// }

// func TestMerklePath_ToHex(t *testing.T) {
// 	t.Parallel()

// 	t.Run("serializes to hex", func(t *testing.T) {
// 		path := MerklePath{
// 			BlockHeight: BRC74JSON.BlockHeight,
// 			Path:        BRC74JSON.Path,
// 		}
// 		hex := path.ToHex()
// 		assert.Equal(t, BRC74Hex, hex)
// 	})
// }

// func TestMerklePath_ComputeRoot(t *testing.T) {
// 	t.Parallel()

// 	t.Run("computes a root", func(t *testing.T) {
// 		path := MerklePath{
// 			BlockHeight: BRC74JSON.BlockHeight,
// 			Path:        BRC74JSON.Path,
// 		}
// 		root := path.ComputeRoot(BRC74TXID1)
// 		assert.Equal(t, BRC74Root, root)
// 	})
// }

// func TestMerklePath_Verify(t *testing.T) {

// 	t.Run("verifies using a ChainTracker", func(t *testing.T) {
// 		path := MerklePath{
// 			BlockHeight: BRC74JSON.BlockHeight,
// 			Path:        BRC74JSON.Path,
// 		}
// 		tracker := &ChainTracker{
// 			IsValidRootForHeight: func(root string, height int) bool {
// 				return root == BRC74Root && height == BRC74JSON.BlockHeight
// 			},
// 		}
// 		result, err := path.Verify(BRC74TXID1, tracker)
// 		assert.NoError(t, err)
// 		assert.True(t, result)
// 	})

// }

// func TestMerklePath_Combine(t *testing.T) {
// 	t.Parallel()

// 	t.Run("combines two paths", func(t *testing.T) {
// 		path0A := make([]PathElement, len(BRC74JSON.Path[0]))
// 		copy(path0A, BRC74JSON.Path[0])
// 		path0B := make([]PathElement, len(BRC74JSON.Path[0]))
// 		copy(path0B, BRC74JSON.Path[0])
// 		pathRest := make([][]PathElement, len(BRC74JSON.Path)-1)
// 		copy(pathRest, BRC74JSON.Path[1:])
// 		path0A = append(path0A[:2], path0A[4:]...)
// 		path0B = path0B[2:]
// 		pathRest = pathRest[1:]
// 		pathRest = append([][]PathElement{path0B}, pathRest...)
// 		pathA := MerklePath{
// 			BlockHeight: BRC74JSON.BlockHeight,
// 			Path:        append([][]PathElement{path0A}, pathRest...),
// 		}
// 		pathB := MerklePath{
// 			BlockHeight: BRC74JSON.BlockHeight,
// 			Path:        append([][]PathElement{path0B}, pathRest...),
// 		}
// 		assert.Equal(t, BRC74Root, pathA.ComputeRoot(BRC74TXID2))
// 		assert.Panics(t, func() { pathA.ComputeRoot(BRC74TXID3) })
// 		assert.Panics(t, func() { pathB.ComputeRoot(BRC74TXID2) })
// 		assert.Equal(t, BRC74Root, pathB.ComputeRoot(BRC74TXID3))
// 		pathA.Combine(pathB)
// 		assert.Equal(t, BRC74Root, pathA.ComputeRoot(BRC74TXID2))
// 		assert.Equal(t, BRC74Root, pathA.ComputeRoot(BRC74TXID3))
// 	})
// }
