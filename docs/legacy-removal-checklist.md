# Legacy Removal Checklist

**Status**: In Progress
**Target**: `internal/app`

## Audit Findings

### Deprecated Fields/Structs

- [x] Scan for `deprecated` comments or tags. (None found in controller/view)
- [x] Removed `Action` field from `ProcessDiaryRequest` (Converted to `CreateDiaryRequest` and `UpdateDiaryRequest`).

### Legacy Logic

- [x] Scan for `data-action` handling. (Found in `DiaryController.ProcessDiary` and `PostController.ProcessPost`)
- [x] Scan for version-based conditional logic. (None found in controllers)

### Legacy Response Formats

- [x] Audit `internal/app/view/response.go` for non-standard responses. (Audit passed, standardized)

### Unused Endpoints

- [x] Identify endpoints not used by Alpine/HTMX. (Removed `/daybook`, `/info`, `/api/post/save`)

## Cleanup Actions

- [x] Remove confirmed legacy code.
- [x] Remove `ProcessDiary` and `ProcessPost` handlers.
- [x] Remove legacy routes in `routes.go`.
- [x] Refactor to RESTful endpoints (`POST /diary`, `PUT /diary/:id`, `POST /post`, `PUT /post/:id`, `DELETE /post/:id`).
