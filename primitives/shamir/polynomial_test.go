package primitives

import (
	"fmt"
	"math/big"
	"testing"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
)

// TestPointInFiniteField verifies the creation and string conversion of points
func TestPointInFiniteField(t *testing.T) {
	x := big.NewInt(10)
	y := big.NewInt(20)

	point := NewPointInFiniteField(x, y)

	if point.X.Cmp(x) != 0 || point.Y.Cmp(y) != 0 {
		t.Errorf("Point creation failed. Expected (%v, %v), got (%v, %v)", x, y, point.X, point.Y)
	}

	str := point.String()
	reconstructedPoint, err := PointFromString(str)
	if err != nil {
		t.Errorf("Failed to reconstruct point from string: %v", err)
	}

	if reconstructedPoint.X.Cmp(point.X) != 0 || reconstructedPoint.Y.Cmp(point.Y) != 0 {
		t.Errorf("Point reconstruction failed. Expected (%v, %v), got (%v, %v)",
			point.X, point.Y, reconstructedPoint.X, reconstructedPoint.Y)
	}
}

// TestPolynomialFromPrivateKey checks if a polynomial is correctly created from a private key
func TestPolynomialFromPrivateKey(t *testing.T) {

	pk, _ := ec.NewPrivateKey()
	threshold := 3

	poly, err := PolynomialFromPrivateKey(pk, threshold)
	if err != nil {
		t.Fatalf("Failed to create polynomial: %v", err)
	}

	if len(poly.Points) != threshold {
		t.Errorf("Incorrect number of points. Expected %d, got %d", threshold, len(poly.Points))
	}

	if poly.Points[0].X.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("First point x-coordinate should be 0, got %v", poly.Points[0].X)
	}

	if poly.Points[0].Y.Cmp(pk.D) != 0 {
		t.Errorf("First point y-coordinate should be the key, got %v", poly.Points[0].Y)
	}

	// Check for uniqueness of x-coordinates
	xCoords := make(map[string]bool)
	for _, point := range poly.Points {
		if xCoords[point.X.String()] {
			t.Errorf("Duplicate x-coordinate found: %v", point.X)
		}
		xCoords[point.X.String()] = true
	}
}

// TestPolynomialValueAt verifies the polynomial evaluation at specific points
func TestPolynomialValueAt(t *testing.T) {
	points := []*PointInFiniteField{
		NewPointInFiniteField(big.NewInt(0), big.NewInt(1)),
		NewPointInFiniteField(big.NewInt(1), big.NewInt(2)),
		NewPointInFiniteField(big.NewInt(2), big.NewInt(3)),
	}
	poly := NewPolynomial(points, 3)

	testCases := []struct {
		x int64
	}{
		{0},
		{1},
		{2},
		{3},
	}

	for _, tc := range testCases {
		x := big.NewInt(tc.x)
		result := poly.ValueAt(x)
		// We're not checking against specific values here, just ensuring it doesn't panic
		t.Logf("Value at x=%d: %v", tc.x, result)
	}
}

func TestFullProcess(t *testing.T) {
	// Create a private key
	privateKey, err := ec.NewPrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	threshold := 3
	totalShares := 5

	// Generate the polynomial
	poly, err := PolynomialFromPrivateKey(privateKey, threshold)
	if err != nil {
		t.Fatalf("Failed to create polynomial: %v", err)
	}

	// Log the generated polynomial points
	t.Logf("Generated polynomial points:")
	for i, point := range poly.Points {
		t.Logf("Point %d: (%v, %v)", i, point.X, point.Y)
	}

	// Generate shares
	shares := make([]*PointInFiniteField, totalShares)
	t.Logf("Generated shares:")
	for i := 0; i < totalShares; i++ {
		x := big.NewInt(int64(i + 1))
		y := poly.ValueAt(x)
		shares[i] = NewPointInFiniteField(x, y)
		t.Logf("Share %d: (%v, %v)", i, shares[i].X, shares[i].Y)
	}

	// Reconstruct the secret using threshold number of shares
	reconstructPoly := NewPolynomial(shares[:threshold], threshold)
	reconstructedSecret := reconstructPoly.ValueAt(big.NewInt(0))

	t.Logf("Original secret: %v", privateKey.D)
	t.Logf("Reconstructed secret: %v", reconstructedSecret)

	if reconstructedSecret.Cmp(privateKey.D) != 0 {
		t.Errorf("Secret reconstruction failed. Expected %v, got %v", privateKey.D, reconstructedSecret)
	}
}

func TestDifferentThresholdsAndShares(t *testing.T) {
	testCases := []struct {
		threshold   int
		totalShares int
	}{
		{2, 3},
		{3, 5},
		{5, 10},
		{10, 20},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Threshold_%d_TotalShares_%d", tc.threshold, tc.totalShares), func(t *testing.T) {
			privateKey, _ := ec.NewPrivateKey()
			poly, err := PolynomialFromPrivateKey(privateKey, tc.threshold)
			if err != nil {
				t.Fatalf("Failed to create polynomial: %v", err)
			}

			shares := make([]*PointInFiniteField, tc.totalShares)
			for i := 0; i < tc.totalShares; i++ {
				x := big.NewInt(int64(i + 1))
				y := poly.ValueAt(x)
				shares[i] = NewPointInFiniteField(x, y)
			}

			reconstructPoly := NewPolynomial(shares[:tc.threshold], tc.threshold)
			reconstructedSecret := reconstructPoly.ValueAt(big.NewInt(0))

			if reconstructedSecret.Cmp(privateKey.D) != 0 {
				t.Errorf("Secret reconstruction failed. Expected %v, got %v", privateKey.D, reconstructedSecret)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	privateKey, _ := ec.NewPrivateKey()

	// Minimum threshold (2)
	t.Run("MinimumThreshold", func(t *testing.T) {
		threshold := 2
		totalShares := 3
		poly, _ := PolynomialFromPrivateKey(privateKey, threshold)
		shares := make([]*PointInFiniteField, totalShares)
		for i := 0; i < totalShares; i++ {
			x := big.NewInt(int64(i + 1))
			y := poly.ValueAt(x)
			shares[i] = NewPointInFiniteField(x, y)
		}
		reconstructPoly := NewPolynomial(shares[:threshold], threshold)
		reconstructedSecret := reconstructPoly.ValueAt(big.NewInt(0))
		if reconstructedSecret.Cmp(privateKey.D) != 0 {
			t.Errorf("Secret reconstruction failed for minimum threshold")
		}
	})

	// Maximum threshold (total shares)
	t.Run("MaximumThreshold", func(t *testing.T) {
		threshold := 10
		totalShares := 10
		poly, _ := PolynomialFromPrivateKey(privateKey, threshold)
		shares := make([]*PointInFiniteField, totalShares)
		for i := 0; i < totalShares; i++ {
			x := big.NewInt(int64(i + 1))
			y := poly.ValueAt(x)
			shares[i] = NewPointInFiniteField(x, y)
		}
		reconstructPoly := NewPolynomial(shares, threshold)
		reconstructedSecret := reconstructPoly.ValueAt(big.NewInt(0))
		if reconstructedSecret.Cmp(privateKey.D) != 0 {
			t.Errorf("Secret reconstruction failed for maximum threshold")
		}
	})
}

func TestReconstructionWithDifferentSubsets(t *testing.T) {
	privateKey, _ := ec.NewPrivateKey()
	threshold := 3
	totalShares := 5

	poly, _ := PolynomialFromPrivateKey(privateKey, threshold)
	shares := make([]*PointInFiniteField, totalShares)
	for i := 0; i < totalShares; i++ {
		x := big.NewInt(int64(i + 1))
		y := poly.ValueAt(x)
		shares[i] = NewPointInFiniteField(x, y)
	}

	subsets := [][]int{
		{0, 1, 2},
		{1, 2, 3},
		{2, 3, 4},
		{0, 2, 4},
	}

	for i, subset := range subsets {
		t.Run(fmt.Sprintf("Subset_%d", i), func(t *testing.T) {
			subsetShares := make([]*PointInFiniteField, threshold)
			for j, idx := range subset {
				subsetShares[j] = shares[idx]
			}
			reconstructPoly := NewPolynomial(subsetShares, threshold)
			reconstructedSecret := reconstructPoly.ValueAt(big.NewInt(0))
			if reconstructedSecret.Cmp(privateKey.D) != 0 {
				t.Errorf("Secret reconstruction failed for subset %v", subset)
			}
		})
	}
}

func TestErrorHandling(t *testing.T) {
	privateKey, _ := ec.NewPrivateKey()

	// Test with invalid threshold (too low)
	_, err := PolynomialFromPrivateKey(privateKey, 1)
	if err == nil {
		t.Errorf("Expected error for threshold too low, got nil")
	}

	// Test with invalid threshold (too high)
	_, err = PolynomialFromPrivateKey(privateKey, 1001)
	if err == nil {
		t.Errorf("Expected error for threshold too high, got nil")
	}

	// Test reconstruction with insufficient shares
	threshold := 3
	poly, _ := PolynomialFromPrivateKey(privateKey, threshold)
	shares := make([]*PointInFiniteField, 2)
	for i := 0; i < 2; i++ {
		x := big.NewInt(int64(i + 1))
		y := poly.ValueAt(x)
		shares[i] = NewPointInFiniteField(x, y)
	}
	reconstructPoly := NewPolynomial(shares, 2)
	reconstructedSecret := reconstructPoly.ValueAt(big.NewInt(0))
	if reconstructedSecret.Cmp(privateKey.D) == 0 {
		t.Errorf("Expected incorrect reconstruction with insufficient shares")
	}
}

func TestConsistency(t *testing.T) {
	privateKey, _ := ec.NewPrivateKey()
	threshold := 3
	totalShares := 5

	for i := 0; i < 10; i++ {
		poly, _ := PolynomialFromPrivateKey(privateKey, threshold)
		shares := make([]*PointInFiniteField, totalShares)
		for j := 0; j < totalShares; j++ {
			x := big.NewInt(int64(j + 1))
			y := poly.ValueAt(x)
			shares[j] = NewPointInFiniteField(x, y)
		}
		reconstructPoly := NewPolynomial(shares[:threshold], threshold)
		reconstructedSecret := reconstructPoly.ValueAt(big.NewInt(0))
		if reconstructedSecret.Cmp(privateKey.D) != 0 {
			t.Errorf("Inconsistent secret reconstruction in run %d", i)
		}
	}
}
