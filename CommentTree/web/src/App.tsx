import { useState } from "react";
import "./App.css";

type SearchItem = {
    id: number;
    parent_id: number | null;
    text: string;
    created_at: string;
    rank: number;
};

type SearchResponse =
    | { status: "ok"; items: SearchItem[]; limit: number; offset: number }
    | { status: "error"; error: string };

type CommentNode = {
    id: number;
    parent_id: number | null;
    text: string;
    created_at: string;
    children?: CommentNode[];
};

type TreeResponse =
    | { status: "ok"; item: CommentNode }
    | { status: "error"; error: string };

type ApiRespOk = { status: "ok"; item?: any };
type ApiRespErr = { status: "error"; error: string };
type ApiResp = ApiRespOk | ApiRespErr;

function TreeLine({
                      node,
                      depth = 0,
                      onReload,
                  }: {
    node: CommentNode;
    depth?: number;
    onReload: () => void;
}) {
    const [openReply, setOpenReply] = useState(false);
    const [replyText, setReplyText] = useState("");
    const [sending, setSending] = useState(false);

    const prefix = depth === 0 ? "" : "↳ ";

    async function submitReply() {
        if (!replyText.trim()) return;

        setSending(true);
        try {
            const r = await fetch("/comments/", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ text: replyText, parent_id: node.id }),
            });

            const data: ApiResp = await r.json();
            if (!r.ok || data.status !== "ok") {
                throw new Error(("error" in data && data.error) ? data.error : `HTTP ${r.status}`);
            }

            setReplyText("");
            setOpenReply(false);
            onReload();
        } catch (e: any) {
            alert(e?.message ?? "post error");
        } finally {
            setSending(false);
        }
    }

    async function deleteThis() {
        const ok = confirm(`Удалить комментарий #${node.id}?`);
        if (!ok) return;

        setSending(true);
        try {
            const r = await fetch(`/comments/${node.id}`, { method: "DELETE" });
            let data: any = null;
            try {
                data = await r.json();
            } catch {
                // если DELETE отдаёт 204 без body — это нормально
            }

            if (!r.ok) {
                const msg = data?.error ? String(data.error) : `HTTP ${r.status}`;
                throw new Error(msg);
            }

            // если удалили корень текущего дерева — дерево исчезнет, но это ок
            onReload();
        } catch (e: any) {
            alert(e?.message ?? "delete error");
        } finally {
            setSending(false);
        }
    }

    return (
        <div className="treeLine" style={{ marginLeft: depth * 18 }}>
            <div className="treeRow">
                <span className="treePrefix">{prefix}</span>
                <span className={depth === 0 ? "treeRootText" : "treeText"}>{node.text}</span>
                <span className="treeMeta">#{node.id}</span>

                <button className="linkBtn" onClick={() => setOpenReply((v) => !v)}>
                    Ответить
                </button>
                <button className="dangerBtn" onClick={deleteThis} disabled={sending}>
                    Удалить
                </button>
            </div>

            {openReply ? (
                <div className="replyBox">
                    <input
                        className="input"
                        value={replyText}
                        onChange={(e) => setReplyText(e.target.value)}
                        placeholder={`Ответ для #${node.id}`}
                    />
                    <button className="button" onClick={submitReply} disabled={sending || !replyText.trim()}>
                        Отправить
                    </button>
                </div>
            ) : null}

            {node.children?.length
                ? node.children.map((ch) => (
                    <TreeLine key={ch.id} node={ch} depth={depth + 1} onReload={onReload} />
                ))
                : null}
        </div>
    );
}

export default function App() {
    const [q, setQ] = useState("Глава");
    const [loading, setLoading] = useState(false);
    const [err, setErr] = useState<string | null>(null);

    const [results, setResults] = useState<SearchItem[]>([]);
    const [tree, setTree] = useState<CommentNode | null>(null);
    const [selectedID, setSelectedID] = useState<number | null>(null);

    // Новый корневой комментарий (без parent_id)
    const [newText, setNewText] = useState("");
    const [creating, setCreating] = useState(false);

    async function loadTree(id: number) {
        setSelectedID(id);
        setTree(null);

        const r = await fetch(`/comments/tree?id=${id}`);
        const data: TreeResponse = await r.json();
        if (!r.ok || data.status !== "ok") {
            throw new Error(("error" in data && data.error) ? data.error : `HTTP ${r.status}`);
        }
        setTree(data.item);
    }

    async function search() {
        setLoading(true);
        setErr(null);
        setResults([]);
        setTree(null);
        setSelectedID(null);

        try {
            const r = await fetch(`/comments/search?q=${encodeURIComponent(q)}&limit=20&offset=0&sort=rank`);
            const data: SearchResponse = await r.json();
            if (!r.ok || data.status !== "ok") {
                throw new Error(("error" in data && data.error) ? data.error : `HTTP ${r.status}`);
            }

            setResults(data.items);

            if (data.items.length > 0) {
                await loadTree(data.items[0].id);
            } else {
                setErr("Ничего не найдено");
            }
        } catch (e: any) {
            setErr(e?.message ?? "unknown error");
        } finally {
            setLoading(false);
        }
    }

    async function createRoot() {
        if (!newText.trim()) return;

        setCreating(true);
        setErr(null);
        try {
            const r = await fetch("/comments/", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ text: newText, parent_id: null }),
            });

            const data: ApiResp = await r.json();
            if (!r.ok || data.status !== "ok") {
                throw new Error(("error" in data && data.error) ? data.error : `HTTP ${r.status}`);
            }

            setNewText("");
            // самый простой UX: после создания делаем поиск по тексту и откроем дерево первого совпадения
            setQ(newText);
            await (async () => {
                const r2 = await fetch(`/comments/search?q=${encodeURIComponent(newText)}&limit=20&offset=0&sort=created_at`);
                const d2: SearchResponse = await r2.json();
                if (r2.ok && d2.status === "ok") {
                    setResults(d2.items);
                    if (d2.items.length > 0) await loadTree(d2.items[0].id);
                }
            })();
        } catch (e: any) {
            setErr(e?.message ?? "create error");
        } finally {
            setCreating(false);
        }
    }

    return (
        <div className="page">
            <header className="header">
                <div className="container">
                    <h1 className="title">Комментарии</h1>

                    <div className="searchRow">
                        <input className="input" value={q} onChange={(e) => setQ(e.target.value)} placeholder="Поиск..." />
                        <button className="button" onClick={search} disabled={loading || !q.trim()}>
                            Найти
                        </button>
                    </div>

                    <div className="createRow">
                        <input
                            className="input"
                            value={newText}
                            onChange={(e) => setNewText(e.target.value)}
                            placeholder="Новый комментарий (корень)..."
                        />
                        <button className="button" onClick={createRoot} disabled={creating || !newText.trim()}>
                            Создать
                        </button>
                    </div>

                    {err ? <div className="error">{err}</div> : null}
                </div>
            </header>

            <main className="main">
                <div className="container grid">
                    <section className="panel">
                        <div className="panelTitle">Совпадения</div>

                        {results.length === 0 ? (
                            <div className="muted">Пока пусто</div>
                        ) : (
                            <div className="list">
                                {results.map((it) => (
                                    <button
                                        key={it.id}
                                        className={"listItem " + (selectedID === it.id ? "active" : "")}
                                        onClick={() => loadTree(it.id)}
                                    >
                                        <div className="listText">{it.text}</div>
                                        <div className="muted">#{it.id}</div>
                                    </button>
                                ))}
                            </div>
                        )}
                    </section>

                    <section className="panel">
                        <div className="panelTitle">Дерево</div>

                        {loading ? <div className="muted">Загрузка...</div> : null}

                        {tree ? (
                            <div className="treeBox">
                                <TreeLine node={tree} onReload={() => loadTree(tree.id)} />
                            </div>
                        ) : (
                            <div className="muted">Сначала выполните поиск</div>
                        )}
                    </section>
                </div>
            </main>
        </div>
    );
}
