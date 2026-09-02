import { useEffect, useRef, useState } from "react";
import { ArrowRight, LogOut, MessageCircle, Send } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Field,
  FieldDescription,
  FieldLabel,
} from "@/components/ui/field";

// server wraps join/leave as "...", notice them in yellow (same as tview client)
function isNotice(m: string) {
  return m.endsWith(" joined the chat...") || m.endsWith(" left the chat...");
}

const timeNow = () =>
  new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });

export default function App() {
  // username input
  const [username, setUsername] = useState("");
  // null hasnt entered yet
  const [name, setName] = useState<string | null>(null);

  // socket stores connection outside use effect to send messages later
  const socket = useRef<WebSocket | null>(null);

  const [text, setText] = useState("");

  const [messages, setMessages] = useState<{ body: string; time: string }[]>([]);

  const endRef = useRef<HTMLDivElement>(null);
  useEffect(() => {
    endRef.current?.scrollIntoView();
  }, [messages]);

  function handleSend(e: { preventDefault: () => void }) {
    e.preventDefault(); // stops the page from reloading
    if (!text.trim()) return; // ignore empty
    socket.current?.send(text);
    setMessages((prev) => [...prev, { body: name + ": " + text, time: timeNow() }]); // local echo (no broadcast to self)
    setText("");
  }

  // sends "exit" so the server broadcasts a clean "left the chat..." then closes
  function handleLeave() {
    socket.current?.send("exit");
    socket.current?.close();
    socket.current = null;
    setText("");
    setMessages([]);
    setUsername("");
    setName(null);
  }

  useEffect(() => {
    if (!name) return;
    // ws = active connection created in this effect
    const ws = new WebSocket("ws://localhost:5173/ws");

    socket.current = ws;

    ws.onopen = () => ws.send(name); // first send is the username

    ws.onmessage = (e) =>
      setMessages((prev) => [
        ...prev,
        { body: e.data as string, time: timeNow() },
      ]);

    return () => {
      ws.close(); // closes the connection when the effect re-runs (StrictMode)
      socket.current = null;
    };
  }, [name]);

  if (!name) {
    return (
      <div className="flex h-dvh items-center justify-center p-6">
        <form
          className="w-full max-w-sm"
          onSubmit={(e) => {
            e.preventDefault();
            setName(username);
          }}
        >
          <Field>
            <FieldLabel htmlFor="input-field-username">Username</FieldLabel>
            <div className="flex items-start gap-2">
              <Input
                id="input-field-username"
                type="text"
                placeholder="Enter your username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
              />
              <Button
                type="submit"
                size="icon"
                aria-label="Join"
              >
                <ArrowRight size={16} />
              </Button>
            </div>
            <FieldDescription>
              Choose a unique username for your chat session.
            </FieldDescription>
          </Field>
        </form>
      </div>
    );
  }

  return (
    <div className="flex h-dvh flex-col">
      <header className="flex items-center justify-between px-5 py-3">
        <div className="flex items-center gap-3">
          <MessageCircle size={18} className="text-foreground/50" />
          <div>
            <h1 className="text-lg font-semibold leading-tight tracking-tight">
              go-chat
            </h1>
            <p className="text-xs leading-tight text-muted-foreground">
              {name}
            </p>
          </div>
        </div>
        <Button
          variant="ghost"
          size="sm"
          onClick={handleLeave}
          className="group text-red-400/90 transition-colors hover:text-red-400 active:text-red-300"
        >
          <LogOut className="mr-1 h-4 w-4 transition-transform duration-200 group-hover:-translate-y-0.5 group-hover:translate-x-0.5 active:translate-y-0" />
          Leave
        </Button>
      </header>

      <main className="flex-1 space-y-1.5 overflow-y-auto px-5 pb-3">
        {messages.map((m, i) => {
          const [user, ...rest] = m.body.split(": ");
          const hasUser = rest.length > 0;
          return (
            <div key={i} className="flex items-baseline gap-2">
              <p
                className={isNotice(m.body) ? "text-amber-400" : "text-foreground/90"}
              >
                {hasUser ? (
                  <>
                    <span className="font-semibold">{user}</span>: {rest.join(": ")}
                  </>
                ) : (
                  m.body
                )}
              </p>
              <span className="ml-auto shrink-0 text-xs tabular-nums text-muted-foreground/80">
                {m.time}
              </span>
            </div>
          );
        })}
        <div ref={endRef} />
      </main>

      <form
        onSubmit={handleSend}
        className="flex items-center gap-2 border-t border-border/60 bg-background/70 px-5 py-3 backdrop-blur-xl"
      >
        <Input
          value={text}
          onChange={(e) => setText(e.target.value)}
          placeholder="Type a message ..."
          className="flex-1"
        />
        <Button type="submit" size="icon" className="shrink-0" aria-label="Send">
          <Send size={16} />
        </Button>
      </form>
    </div>
  );
}