# /seo-run — SEO 파이프라인 실행 + 결과 분석

## 설명
6개 전문 에이전트 파이프라인을 실행하여 종합 SEO 캠페인을 생성하고 분석 리포트를 제공한다.
Claude Code가 각 에이전트의 시스템 프롬프트를 참고하여 직접 결과를 생성한다 (외부 API 키 불필요).

## 사용법
```
/seo-run <SEO 캠페인 요청>
```
예시: `/seo-run 서울 강남 카페 SEO 최적화`

## 실행 절차

### Step 1: 에이전트 프롬프트 로드
다음 6개 파일을 모두 읽는다:
- `prompts/knowledge_guide.txt`
- `prompts/planner.txt`
- `prompts/developer.txt`
- `prompts/executor.txt`
- `prompts/trust_manager.txt`
- `prompts/analyst.txt`

### Step 2: 파이프라인 실행 (DAG 순서)
사용자의 `$ARGUMENTS`를 초기 요청으로 하여, 각 에이전트의 시스템 프롬프트를 역할 지침으로 삼아 순차적으로 결과를 생성한다.

**Layer 1 — knowledge_guide:**
- 시스템 프롬프트: `prompts/knowledge_guide.txt`
- 입력: 사용자 요청 (`$ARGUMENTS`)
- 출력: 키워드 리서치, 타겟 분석, 전략 수립
- 결과를 `=== [knowledge_guide] START ===` / `=== [knowledge_guide] END ===` 구분자로 출력

**Layer 2 — planner:**
- 시스템 프롬프트: `prompts/planner.txt`
- 입력: knowledge_guide 결과
- 출력: 콘텐츠 아웃라인, 실행 계획
- 결과를 구분자로 출력

**Layer 3 — developer + executor:**
- **developer**: `prompts/developer.txt` 참고, planner 결과 기반 → 메타태그, JSON-LD, 사이트맵
- **executor**: `prompts/executor.txt` 참고, planner 결과 기반 → 블로그 글, 소셜미디어 카피
- 각각 구분자로 출력

**Layer 4 — trust_manager:**
- 시스템 프롬프트: `prompts/trust_manager.txt`
- 입력: developer + executor 결과
- 출력: E-E-A-T 점수, 수정 요청
- **재시도 로직**: E-E-A-T 평균 점수가 6 미만이면 Layer 3(developer+executor) 재실행 (최대 2회)
- 결과를 구분자로 출력

**Layer 5 — analyst:**
- 시스템 프롬프트: `prompts/analyst.txt`
- 입력: 전체 에이전트 결과
- 출력: KPI 예측, 최적화 제안
- 결과를 구분자로 출력

### Step 3: 에이전트별 품질 평가
각 에이전트의 출력을 다음 기준으로 평가한다:

**knowledge_guide**:
- 키워드 수 (최소 5개 이상인가?)
- 검색 의도 분류가 되어있는가?
- 한국 시장 특화 키워드 포함 여부
- 브랜드 포지셔닝 스테이트먼트 완성도 ("We are X for Y who Z" 형태)
- USP 명확성 및 차별점 3개 도출 여부
- 경쟁 브랜드 포지셔닝 맵 포함 여부

**planner**:
- 콘텐츠 아웃라인의 구체성
- 실행 가능한 작업 계획인가?
- 브랜드 보이스 가이드라인 완성도 (톤, 스타일 규칙, 선호/금지 용어)
- 섹션별 브랜드 메시지 매핑 여부
- 플랫폼별 보이스 변주 정의 여부

**developer**:
- 메타태그 완성도 (title, description, OG tags)
- JSON-LD 스키마 정확성
- 기술적 SEO 체크리스트 충실도
- Organization JSON-LD 포함 여부 (name, logo, sameAs)
- 브랜드 메타태그 (og:site_name, theme-color) 설정 여부

**executor**:
- 블로그 글 길이와 구조 (H2/H3 사용)
- 키워드 자연스러운 포함
- 소셜미디어 카피 플랫폼별 적합성
- 브랜드 보이스 적용 여부 (brand_voice_applied)
- USP 멘션 횟수 (최소 2회)
- 브랜드 해시태그 포함 여부
- 브랜드 스토리텔링 요소 포함 여부

**trust_manager**:
- E-E-A-T 각 항목 점수
- 수정 요청이 있었는지, 몇 회 재시도했는지
- Brand Authenticity 점수 (VoiceConsistency, USPClarity, PositioningAlignment)
- 브랜드 일관성 체크 결과 (voice_tone, message_hierarchy, competitor_terms)

**analyst**:
- KPI 예측의 현실성
- 개선 제안의 구체성과 실행 가능성
- 브랜드 KPI 포함 여부 (인지도 예측, 포지셔닝 명확성, 보이스 일관성, 차별화 수준)
- 경쟁 브랜드 대비 분석 포함 여부

### Step 4: 종합 리포트 출력
다음 형식으로 분석 결과를 제공한다:

```
## SEO 캠페인 실행 결과

### 캠페인 점수: X/10

### 브랜딩 품질 요약
| 항목 | 점수 | 비고 |
|------|------|------|
| 포지셔닝 명확성 | X/10 | ... |
| 보이스 일관성 | X/10 | ... |
| USP 전달력 | X/10 | ... |
| 브랜드 진정성 | X/10 | ... |
| 차별화 수준 | high/medium/low | ... |

### 에이전트별 평가
| 에이전트 | 점수 | 핵심 산출물 | 비고 |
|---------|------|-----------|------|
| knowledge_guide | X/10 | 키워드 N개 | ... |
| planner | X/10 | 아웃라인 N개 섹션 | ... |
| developer | X/10 | 메타태그+JSON-LD | ... |
| executor | X/10 | 블로그 N자 | ... |
| trust_manager | X/10 | E-E-A-T 평균 X점 | ... |
| analyst | X/10 | KPI N개 | ... |

### 강점
- ...

### 약점 / 개선 필요
- ...

### 다음 단계 제안
1. ...
```

## 에러 처리
- 프롬프트 파일 누락 시: 해당 에이전트를 기본 역할 지식으로 실행하고 경고 출력
- 결과 품질 미달 시: trust_manager 재시도 로직 적용 (최대 2회)
