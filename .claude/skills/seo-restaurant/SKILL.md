# /seo-restaurant — 음식점/카페 전문 SEO 캠페인

## 설명
음식점/카페에 특화된 SEO 캠페인을 실행한다. 도메인 지식을 파이프라인 요청에 주입하고 결과를 음식점 맥락에서 평가한다.
Claude Code가 직접 에이전트 역할을 수행한다 (외부 API 키 불필요).

## 사용법
```
/seo-restaurant <업소명> <지역> [업종]
```
예시: `/seo-restaurant 맛있는집 서울 강남 한식당`

## 실행 절차

### Step 1: 도메인 지식 로드
`.claude/skills/seo-restaurant/domain-knowledge.md`를 읽어 음식점 SEO 전문 지식을 로드한다.

### Step 2: 에이전트 프롬프트 로드
다음 6개 파일을 모두 읽는다:
- `prompts/knowledge_guide.txt`
- `prompts/planner.txt`
- `prompts/developer.txt`
- `prompts/executor.txt`
- `prompts/trust_manager.txt`
- `prompts/analyst.txt`

### Step 3: 강화된 요청 구성
사용자의 `$ARGUMENTS`에서 업소명, 지역, 업종을 추출하고, 도메인 지식을 결합:

```
[음식점/카페 SEO 캠페인]
업소명: {업소명}
지역: {지역}
업종: {업종}

핵심 목표:
1. 지역 맛집 검색 "{지역} {업종} 맛집" 상위 노출
2. 네이버 Place 최적화
3. 배달앱(배민/쿠팡이츠) 연동 SEO
4. Restaurant JSON-LD 구조화 데이터
5. Instagram 비주얼 콘텐츠 전략
6. 음식점 브랜드 포지셔닝 확립

추가 컨텍스트:
- 메뉴, 가격대, 분위기, 특징을 콘텐츠에 반영
- 리뷰/평점 관리 전략 포함
- 시즌 메뉴/이벤트 키워드 대응
- 사진 중심 소셜미디어 전략

브랜드 포지셔닝:
- 음식점 고유의 USP 도출 (비법 레시피, 분위기, 서비스 차별점 등)
- 브랜드 성격 정의 (정통적/트렌디/가성비/프리미엄 등)
- 셰프/오너의 철학을 브랜드 스토리에 반영
- 맛집 브랜드 보이스 설정 (맛 표현 어휘, 감성적 문체 등)
```

### Step 4: 파이프라인 실행 (DAG 순서)
강화된 요청을 초기 입력으로 하여, 각 에이전트의 시스템 프롬프트를 역할 지침으로 삼아 순차적으로 결과를 생성한다.
도메인 지식(domain-knowledge.md)의 키워드, JSON-LD 스키마, 네이버 Place 전략을 각 에이전트 실행 시 참고한다.

**Layer 1 — knowledge_guide:**
- 입력: 강화된 요청
- 출력: 맛집 키워드 리서치, 타겟 분석, 전략 수립
- 도메인 지식의 키워드 카테고리(지역 맛집, 업종별, 상황별, 카페 특화) 참고

**Layer 2 — planner:**
- 입력: knowledge_guide 결과
- 출력: 콘텐츠 아웃라인, 실행 계획

**Layer 3 — developer + executor:**
- **developer**: 메타태그, Restaurant JSON-LD 스키마, 사이트맵
- **executor**: 블로그 글, 소셜미디어 카피, Instagram 콘텐츠

**Layer 4 — trust_manager:**
- E-E-A-T 평가 (점수 < 6이면 Layer 3 재실행, 최대 2회)

**Layer 5 — analyst:**
- KPI 예측, 최적화 제안

각 에이전트 결과는 `=== [agent_name] START/END ===` 구분자로 출력한다.

### Step 5: 음식점 도메인 특화 평가
일반 파이프라인 평가에 추가로:

- **Restaurant JSON-LD**: `@type: Restaurant` 스키마 포함 여부
- **네이버 Place**: 업체 정보 구조화, 리뷰 관리 전략
- **배달앱 최적화**: 배달의민족/쿠팡이츠 키워드
- **Instagram 전략**: 음식 사진, 해시태그, 릴스 전략
- **계절 메뉴**: 시즌별 키워드 대응
- **리뷰 관리**: 별점/리뷰 획득 전략

### Step 6: 리포트 출력
```
## 음식점 SEO 캠페인 결과: {업소명}

### 캠페인 점수: X/10

### 음식점 특화 체크리스트
- [ ] Restaurant JSON-LD 스키마
- [ ] 네이버 Place 최적화
- [ ] 배달앱 키워드 전략
- [ ] Instagram 콘텐츠 플랜
- [ ] 메뉴 구조화 데이터
- [ ] 리뷰 관리 전략
- [ ] 음식점 브랜드 포지셔닝 (USP, 차별점)
- [ ] 브랜드 보이스 가이드라인 (맛 표현, 감성 톤)
- [ ] Organization JSON-LD (업소명, 로고, 소셜 프로필)

### 에이전트별 평가
(표 형식)

### 음식점 SEO 개선 제안
1. ...
```
