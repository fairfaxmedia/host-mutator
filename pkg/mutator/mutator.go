package mutator

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	admission "k8s.io/api/admission/v1beta1"
	networking "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var baseDomain string

func init() {
	baseDomain = os.Getenv("BASE_DOMAIN")
}

// Mutate receives an AdmissionReview (request), adds a response, and then returns it
// Its goal is to append a base domain to the host value in a kubernetes ingress resource
func Mutate(body []byte) ([]byte, error) {

	// unmarshal the request
	admReview := admission.AdmissionReview{}
	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, fmt.Errorf("unmarshaling request failed with %s", err)
	}
	request := admReview.Request

	// given an empty request, return an empty response
	if request == nil {
		return []byte{}, nil
	}

	// get the ingress resource from the request
	var ingress *networking.Ingress
	if err := json.Unmarshal(request.Object.Raw, &ingress); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ingress json object %v", err)
	}

	// set the response options
	response := &admission.AdmissionResponse{}
	response.Allowed = true
	response.UID = request.UID
	patchType := admission.PatchTypeJSONPatch
	response.PatchType = &patchType
	response.AuditAnnotations = map[string]string{
		"mutated-host": "true",
	}

	// build a JSONPatch for each host rule
	patches := []map[string]string{}
	for i, rule := range ingress.Spec.Rules {
		patch := map[string]string{
			"op":    "replace",
			"path":  fmt.Sprintf("/spec/rules/%d/host", i),
			"value": strings.Join([]string{rule.Host, baseDomain}, "."),
		}
		patches = append(patches, patch)
	}

	// add the patches to the response
	jsonPatches, err := json.Marshal(patches)
	if err != nil {
		return nil, err
	}
	response.Patch = jsonPatches

	// set the result as success
	response.Result = &metav1.Status{
		Status: "Success",
	}

	// Add the response to the AdmissionReview and return it
	admReview.Response = response
	responseBody, err := json.Marshal(admReview)
	if err != nil {
		return nil, err
	}
	return responseBody, nil

}
