# CommentTree

Сервис комментариев с деревом ответов, поиском, созданием и удалением + простой Web UI.

## Запуск backend
```bash
go run ./cmd/app
```
Backend: http://localhost:8080

Запуск web
```azure
cd web
npm install
npm run dev
```
Frontend: http://localhost:5173/

API

– POST /comments — создание комментария (с указанием родительского);

– GET /comments?parent={id} — получение комментария и всех вложенных;

– DELETE /comments/{id} — удаление комментария и всех вложенных под ним.
