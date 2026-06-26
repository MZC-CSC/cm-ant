// Package securitycases 는 보안 게이트(gosec) 발화 검증용 더미 코드다.
// 의도적으로 약한 패턴을 담았으며 실제 코드에서 사용하지 않는다. (테스트 전용 · 미머지)
package securitycases

import (
	"fmt"
	"math/rand"
)

// apiSecret 는 하드코딩 시크릿 — gosec G101 유발용 더미 값.
// (저엔트로피 값은 gosec 엔트로피 임계값에 안 걸리므로 고엔트로피 더미 사용)
const apiSecret = "kJ8sJ2mNp4Qr7Tv0WxZ1aBc3DeF5gHi"

// WeakToken 은 암호학적으로 안전하지 않은 난수 사용 — gosec G404 유발용.
func WeakToken() string {
	return fmt.Sprintf("tok-%d", rand.Int())
}

// DummyUse 는 위 더미 값들을 참조해 컴파일이 되도록 한다(미사용 경고 방지).
func DummyUse() string {
	return apiSecret + ":" + WeakToken()
}
