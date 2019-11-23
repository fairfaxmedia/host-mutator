package mutator

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	v1beta1 "k8s.io/api/admission/v1beta1"
)

func TestMutateJSON(t *testing.T) {

	if err := os.Setenv("BASE_DOMAIN", "test.run"); err != nil {
		t.Fatal(err)
	}

	admissionReviewJSON, err := ioutil.ReadFile("testdata/admission-review.json")
	if err != nil {
		t.Fatal(err)
	}
	response, err := Mutate(admissionReviewJSON)
	if err != nil {
		t.Errorf("Failed to mutate AdmissionRequest %s with error %s", string(response), err)
	}

	r := v1beta1.AdmissionReview{}
	err = json.Unmarshal(response, &r)
	assert.NoError(t, err, "failed to unmarshal with error %s", err)

	rr := r.Response
	assert.Equal(t, `[{"op":"replace","path":"/spec/rules/0/host","value":"kubernetes-dashboard.test.run"}]`, string(rr.Patch))
	assert.Contains(t, rr.AuditAnnotations, "mutated-host")

}
