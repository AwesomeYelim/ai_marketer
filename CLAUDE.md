# AI Marketer — 멀티 에이전트 SEO 시스템

## 프로젝트 개요
Go 기반 멀티 에이전트 SEO 파이프라인. 6개 전문 에이전트가 DAG 구조로 협업하여 종합 SEO 캠페인을 자동 생성한다.
Claude Code가 **상위 에이전트**로서 `prompts/*.txt`를 읽고 직접 파이프라인을 실행한다 (외부 API 키 불필요).

## 아키텍처

```
Claude Code (/seo-run, /seo-church, /seo-restaurant, /seo-fashion)
    │
    ├── CLAUDE.md                    ← 프로젝트 컨텍스트 (항상 로드)
    ├── .claude/skills/              ← 슬래시 커맨드 정의
    ├── prompts/*.txt                ← 에이전트 시스템 프롬프트 (Claude Code가 직접 참고)
    │
    └── Go Pipeline (선택적, ANTHROPIC_API_KEY 필요 시)
         ├── go run main.go run "..."
         └── stdout: JSON 결과
```

**핵심 원리**: Claude Code가 각 에이전트의 시스템 프롬프트를 읽고 해당 역할을 수행한다. Go 파이프라인은 독립 실행용으로 유지한다.

## 기술 스택
- **언어**: Go 1.23
- **LLM**: Claude Sonnet 4.6 (Anthropic SDK)
- **CLI**: Cobra
- **설정**: YAML (`config.yaml`)
- **상위 에이전트**: Claude Code (슬래시 커맨드)

## 실행 방법

### Claude Code 스킬 (권장, API 키 불필요)
```
/seo-run 서울 강남 카페 SEO 최적화
/seo-church 사랑의교회 서울 서초구
/seo-restaurant 맛있는집 서울 강남 한식당
/seo-fashion 모다브랜드 20대여성 미니멀
```

### Go CLI 직접 실행 (ANTHROPIC_API_KEY 필요)
```bash
go build -o ai-marketer ./...
go run main.go run "서울 강남 카페 SEO 최적화"
go run main.go -c config.yaml run "요청 내용"
```

## 파이프라인 아키텍처 (6 에이전트 DAG)

```
Layer 1: knowledge_guide        ← 키워드 리서치, 전략 수립, 브랜드 포지셔닝
           ↓
Layer 2: planner                ← 실행 계획, 콘텐츠 아웃라인, 브랜드 보이스 설계
           ↓
Layer 3: developer ∥ executor   ← 기술 SEO + 브랜드 스키마 ∥ 콘텐츠 생성 + 브랜드 보이스 적용 (병렬)
           ↓
Layer 4: trust_manager          ← E-E-A-T 평가 + 브랜드 진정성 평가, 품질 관리
           ↓                       (E-E-A-T 또는 브랜드 점수 < 6이면 Layer 3 재실행, 최대 2회)
Layer 5: analyst                ← 성과 예측, 최적화 제안, 브랜드 KPI 분석
```

## 에이전트 상세

| 에이전트 | 프롬프트 | Go 코드 | 입력 | 출력 |
|---------|---------|---------|------|------|
| knowledge_guide | `prompts/knowledge_guide.txt` | `internal/agent/knowledge_guide.go` | 사용자 요청 | 키워드, 타겟 분석, 전략, **브랜드 포지셔닝** |
| planner | `prompts/planner.txt` | `internal/agent/planner.go` | knowledge_guide 결과 | 콘텐츠 아웃라인, 작업 계획, **브랜드 보이스 가이드라인** |
| developer | `prompts/developer.txt` | `internal/agent/developer.go` | planner 결과 | 메타태그, JSON-LD, 사이트맵, **Organization 스키마, 브랜드 메타태그** |
| executor | `prompts/executor.txt` | `internal/agent/executor.go` | planner 결과 | 블로그 글, 소셜미디어 카피, **USP 삽입, 브랜드 해시태그** |
| trust_manager | `prompts/trust_manager.txt` | `internal/agent/trust_manager.go` | developer + executor 결과 | E-E-A-T 점수, 수정 요청, **브랜드 진정성 점수** |
| analyst | `prompts/analyst.txt` | `internal/agent/analyst.go` | 전체 결과 | KPI 예측, 개선 제안, **브랜드 KPI** |

## 코드 구조

```
ai_marketer/
├── main.go                         # CLI 진입점 (Cobra)
├── config.yaml                     # LLM/파이프라인 설정
├── prompts/                        # 에이전트 시스템 프롬프트 (수정 가능)
│   ├── knowledge_guide.txt
│   ├── planner.txt
│   ├── developer.txt
│   ├── executor.txt
│   ├── trust_manager.txt
│   └── analyst.txt
├── internal/
│   ├── agent/                      # 6개 에이전트 구현
│   │   ├── agent.go                # Agent 인터페이스, BaseAgent
│   │   ├── registry.go             # 에이전트 레지스트리
│   │   └── [agent_name].go         # 개별 에이전트
│   ├── orchestrator/               # 파이프라인 실행 엔진
│   │   ├── orchestrator.go         # 오케스트레이터
│   │   └── pipeline.go             # 파이프라인 정의
│   ├── message/                    # 메시지/컨텍스트
│   │   ├── message.go              # Message, AgentResponse
│   │   └── context.go              # WorkflowContext
│   ├── config/                     # 설정 관리
│   └── llm/                        # Anthropic SDK 래퍼
│       ├── client.go               # Complete(), CompleteStructured[T]()
│       └── schema.go               # JSON 스키마 생성
└── pkg/types/types.go              # 공유 타입
```

## 출력 포맷
파이프라인 실행 시 각 에이전트 결과가 구분자와 함께 출력된다:
```
=== [agent_name] START ===
{ JSON 결과 }
=== [agent_name] END ===
```

## 브랜딩 레이어

SEO + 브랜딩이 융합된 캠페인을 생성하도록 각 에이전트에 브랜딩 역할이 추가되어 있다.

| 에이전트 | 브랜딩 역할 | 핵심 구조체 |
|---------|-----------|-----------|
| knowledge_guide | 브랜드 전략가 — 포지셔닝, USP, 경쟁 맵 | `BrandPositioning`, `CompetitorBrand` |
| planner | 보이스 설계자 — 톤, 문체, 플랫폼별 변주 | `BrandVoiceGuideline`, `SectionMessage` |
| developer | 스키마 전문가 — Organization JSON-LD, 브랜드 메타 | `OrganizationSchema`, `BrandMetaTags` |
| executor | 보이스 카피라이터 — USP 삽입, 스토리텔링 | `BrandVoiceApplied`, `BrandHashtags` |
| trust_manager | 진정성 검증자 — 보이스 일관성, USP 명확성 | `BrandAuthenticityScore`, `BrandConsistencyCheck` |
| analyst | 브랜드 성과 분석 — 인지도, 차별화, 일관성 | `BrandKPIs` |

**브랜드 품질 게이트**: trust_manager의 Brand Authenticity 점수가 6 미만이면 E-E-A-T와 마찬가지로 Layer 3 재실행.

## 타겟 산업
- **교회**: 예배 안내, 새신자 유입, 지역 SEO, Church JSON-LD, 교회 브랜드 아이덴티티
- **음식점/카페**: 맛집 키워드, 배달앱, Naver Place, Restaurant JSON-LD, 음식점 브랜드 포지셔닝
- **의류/브랜드**: 코디 키워드, 무신사/에이블리, Product JSON-LD, 시즌 룩북, 패션 브랜드 아이덴티티

## Claude Code 스킬 (슬래시 커맨드)
- `/seo-run` — 파이프라인 실행 + 결과 분석
- `/seo-review` — 파이프라인 출력 심층 리뷰
- `/seo-improve` — 에이전트 프롬프트 최적화
- `/seo-prompt` — 단일 에이전트 프롬프트 심층 분석
- `/seo-church` — 교회 전문 SEO 캠페인
- `/seo-restaurant` — 음식점/카페 전문 SEO 캠페인
- `/seo-fashion` — 의류/브랜드 전문 SEO 캠페인

## 주요 규칙
- 프롬프트 수정 시 반드시 사용자 확인 후 적용
- 한국어 콘텐츠 중심 (Naver/Kakao 검색 최적화 포함)
- E-E-A-T 점수 6 미만 시 자동 재시도 (최대 2회)
- Brand Authenticity 점수 6 미만 시 자동 재시도 (최대 2회)
- 모든 콘텐츠에 브랜드 보이스 가이드라인 일관 적용 필수
