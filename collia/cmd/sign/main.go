// collia/cmd/sign 是发布方离线签名工具：生成 Ed25519 密钥对，并为升级包
// 生成 .sig 侧车文件（对二进制 SHA-256 摘要签名）。私钥只在发布环境使用，
// 不进入 Beacon 服务端与 Agent 镜像；公钥写入服务端 config.toml 的
// [AgentUpdate].PublicKey 并随安装 config 下发到 Agent 固定。
package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	gen := flag.String("gen", "", "生成新密钥对：私钥写入该路径（base64，0600），公钥打印到标准输出")
	key := flag.String("key", "", "签名所用私钥文件（base64）")
	flag.Parse()

	switch {
	case *gen != "":
		if err := genKey(*gen); err != nil {
			fail(err)
		}
	case *key != "" && flag.NArg() == 1:
		if err := signFile(*key, flag.Arg(0)); err != nil {
			fail(err)
		}
	default:
		fmt.Fprintln(os.Stderr, "用法: sign -gen <private-key-path> | sign -key <private-key-path> <file>")
		fmt.Fprintln(os.Stderr, "  -gen     生成密钥对并打印公钥（写入服务端 [AgentUpdate].PublicKey）")
		fmt.Fprintln(os.Stderr, "  -key     用私钥为 <file> 生成 <file>.sig 侧车签名文件")
		os.Exit(2)
	}
}

func genKey(path string) error {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(base64.StdEncoding.EncodeToString(priv)+"\n"), 0o600); err != nil {
		return err
	}
	fmt.Printf("私钥已写入: %s（请离线保管，勿进镜像与仓库）\n", path)
	fmt.Printf("公钥（写入服务端 config.toml [AgentUpdate].PublicKey）: %s\n", base64.StdEncoding.EncodeToString(pub))
	return nil
}

func signFile(keyPath, file string) error {
	keyB64, err := os.ReadFile(keyPath)
	if err != nil {
		return err
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(keyB64)))
	if err != nil {
		return fmt.Errorf("decode private key: %w", err)
	}
	if len(raw) != ed25519.PrivateKeySize {
		return fmt.Errorf("私钥应为 %d 字节，得到 %d", ed25519.PrivateKeySize, len(raw))
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	digest := sha256.Sum256(data)
	sig := ed25519.Sign(ed25519.PrivateKey(raw), digest[:])
	out := file + ".sig"
	if err := os.WriteFile(out, []byte(base64.StdEncoding.EncodeToString(sig)+"\n"), 0o644); err != nil {
		return err
	}
	fmt.Printf("已写入签名侧车文件: %s\n", out)
	return nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "sign:", err)
	os.Exit(1)
}
