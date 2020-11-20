package domain

import (
    "testing"
)

func TestPrepareReferencePath(t *testing.T) {
    var expectedPath = []string{
        "./config/service.yaml",
        "./config/service.yaml",
        "./service.yaml",
        "",
        "./service.yaml",
    }
    var expectedElement = []string{
        "service",
        "service/elem",
        "service",
        "service",
        "",
    }
    var refPath = []string{
        "service.yaml#service/",
        "./service.yaml#service/elem",
        "../service.yaml#service",
        "#service",
        "../service.yaml",
    }
    var path, element string

    for i, s := range refPath {
        path, element = NewYamlReader(make(map[string]interface{})).prepareReferencePath(s)
        if expectedPath[i] != path {
            t.Errorf("niewłaściwa ścieżka, jest: %s, a powinna być: %s.", path, expectedPath[i])
        }
        if expectedElement[i] != element {
            t.Errorf("niewłaściwa nazwa elementu, jest: %s, a powinna być: %s.", element, expectedElement[i])
        }
    }
}
