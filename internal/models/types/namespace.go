package types

import "github.com/TrianaLab/awasm-portfolio/internal/models"

// Namespace is cluster-scoped: it has no namespace of its own and cannot
// be owned by another resource.
type Namespace struct {
	models.Meta `json:",inline" yaml:",inline"`
}

func (n *Namespace) GetNamespace() string                      { return "" }
func (n *Namespace) SetNamespace(string)                       {}
func (n *Namespace) SetOwnerReference(_ models.OwnerReference) {}
