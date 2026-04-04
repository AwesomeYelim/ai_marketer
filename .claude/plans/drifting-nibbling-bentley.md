# 활용 가이드 + 문서화 시스템 구축

## Context
SEO 멀티 에이전트 파이프라인과 7개 슬래시 커맨드가 구현 완료되었으나, **어떻게 활용해야 하는지** 안내가 없고 **결과물을 문서로 관리**하는 기능이 없다.

3가지 목표:
1. **사용 가이드** — 팀/동료가 바로 쓸 수 있는 온보딩 + 개인 워크플로우 참고용
2. **고객 제공용 보고서** — 캠페인 결과를 교회/음식점/브랜드 고객에게 전달할 문서 양식
3. **문서화 역할** — 결과 보고서 자동 생성, 프롬프트 변경 이력, 캠페인 아카이브

## 생성 파일 목록 (14개)

### 가이드 문서 (3개)

#### 1. `docs/QUICK_START.md` — 5분 빠른 시작
- 환경 설정 (Go, API Key, Claude Code)
- 첫 캠페인 실행 (`/seo-run`)
- 결과 확인 및 다음 단계

#### 2. `docs/USAGE_GUIDE.md` — 종합 사용 가이드
- 전체 슬래시 커맨드 카탈로그 (7개 + 신규 3개)
- 각 커맨드별 사용 시나리오, 입력/출력 예시
- 커맨드 조합 워크플로우 (실행→리뷰→개선 사이클)
- FAQ / 트러블슈팅

#### 3. `docs/WORKFLOWS.md` — 실전 워크플로우 가이드
- 도메인별 시나리오 (교회/음식점/패션)
- 단계별 워크플로우 다이어그램
- 반복 캠페인 운영 패턴
- 프롬프트 고도화 사이클
- 보고서 생성 및 고객 전달 프로세스

### 고객 보고서 템플릿 (4개)

#### 4. `docs/client-templates/base-report.md` — 범용 SEO 보고서 템플릿
- 캠페인 개요, 키워드 분석, 콘텐츠 전략
- 기술 SEO 체크리스트, E-E-A-T 평가
- KPI 목표, 다음 단계 제안
- 고객이 이해하기 쉬운 비기술적 표현

#### 5. `docs/client-templates/church-report.md` — 교회 전용 보고서 템플릿
- 새신자 유입 전략, 예배 안내 SEO
- Church JSON-LD 적용 가이드
- 네이버 블로그/YouTube 전략

#### 6. `docs/client-templates/restaurant-report.md` — 음식점 전용 보고서 템플릿
- 맛집 키워드 전략, 네이버 Place 최적화
- 배달앱/Instagram 전략
- 리뷰 관리 전략

#### 7. `docs/client-templates/fashion-report.md` — 의류 전용 보고서 템플릿
- 코디/스타일 키워드 전략
- 무신사/에이블리 최적화
- 인플루언서 협업, 시즌 룩북 전략

### 신규 스킬 (3개)

#### 8. `.claude/skills/seo-report/SKILL.md` — `/seo-report` 보고서 자동 생성
- 직전 `/seo-run` 결과를 고객용 보고서로 변환
- 도메인(교회/음식점/패션) 자동 감지 → 해당 템플릿 적용
- `output/reports/` 디렉토리에 마크다운 파일 저장
- 파일명: `{날짜}_{고객명}_SEO_보고서.md`

#### 9. `.claude/skills/seo-archive/SKILL.md` — `/seo-archive` 캠페인 아카이브
- 사용법: `/seo-archive save "설명"` — 현재 캠페인 결과 저장
- 사용법: `/seo-archive list` — 저장된 캠페인 목록
- 사용법: `/seo-archive compare 1 2` — 두 캠페인 비교
- `output/archives/` 디렉토리에 날짜별 JSON 저장

#### 10. `.claude/skills/seo-changelog/SKILL.md` — `/seo-changelog` 프롬프트 변경 이력
- 사용법: `/seo-changelog` — 전체 변경 이력 보기
- 사용법: `/seo-changelog knowledge_guide` — 특정 에이전트 이력
- `/seo-improve`로 프롬프트 수정 시 자동으로 `prompts/CHANGELOG.md`에 기록
- 변경 사유, 이전/이후 요약, 날짜 포함

### 지원 파일 (2개)

#### 11. `prompts/CHANGELOG.md` — 프롬프트 변경 이력 파일
- 초기 버전 기록 (v1.0)
- 변경 시 `/seo-changelog`가 자동 추가

#### 12. `output/.gitkeep` — 출력 디렉토리 생성
- `output/reports/` — 보고서 저장소
- `output/archives/` — 캠페인 아카이브 저장소

### 업데이트 (2개)

#### 13. `CLAUDE.md` — 업데이트
- 신규 스킬 3개 (`/seo-report`, `/seo-archive`, `/seo-changelog`) 추가
- 문서 구조 (`docs/`, `output/`) 추가
- 활용 가이드 참조 링크

#### 14. `.claude/skills/seo-improve/SKILL.md` — 업데이트
- Step 5 추가: 프롬프트 수정 시 `prompts/CHANGELOG.md`에 변경 이력 자동 기록

## 구현 순서
1. `docs/` 디렉토리 구조 생성 + `output/` 디렉토리 생성
2. 가이드 문서 3개 (QUICK_START, USAGE_GUIDE, WORKFLOWS)
3. 고객 보고서 템플릿 4개
4. 신규 스킬 3개 (seo-report, seo-archive, seo-changelog)
5. 지원 파일 (CHANGELOG.md, .gitkeep)
6. 기존 파일 업데이트 (CLAUDE.md, seo-improve)

## 검증
1. `docs/` 디렉토리에 가이드 3개 + 템플릿 4개 확인
2. `/seo-report` → 템플릿 기반 보고서 생성 확인
3. `/seo-archive save "테스트"` → 아카이브 저장 확인
4. `/seo-changelog` → 변경 이력 조회 확인
5. `CLAUDE.md`에 신규 스킬/문서 반영 확인
