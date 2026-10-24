package tunnel

import "testing"

func TestAgentTunnelRegistrationPayloadIncludesMetadata(t *testing.T) {
	tunnel := NewAgentTunnel("127.0.0.1:17000", "agent-1")
	tunnel.SetJoinToken("join-token")
	tunnel.SetMetadata("v3.0.6", "linux", "arm64")

	got := tunnel.registrationPayload()
	if got.AgentID != "agent-1" {
		t.Fatalf("AgentID = %q, want agent-1", got.AgentID)
	}
	if got.JoinToken != "join-token" {
		t.Fatalf("JoinToken = %q, want join-token", got.JoinToken)
	}
	if got.Version != "v3.0.6" || got.OS != "linux" || got.Arch != "arm64" {
		t.Fatalf("metadata = version %q os %q arch %q", got.Version, got.OS, got.Arch)
	}
}
