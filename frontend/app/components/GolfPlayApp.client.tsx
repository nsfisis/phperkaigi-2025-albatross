import type { components } from "../.server/api/schema";
import { useState, useEffect } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { useDebouncedCallback } from "use-debounce";
import GolfPlayAppConnecting from "./GolfPlayApps/GolfPlayAppConnecting";
import GolfPlayAppWaiting from "./GolfPlayApps/GolfPlayAppWaiting";
import GolfPlayAppStarting from "./GolfPlayApps/GolfPlayAppStarting";
import GolfPlayAppGaming from "./GolfPlayApps/GolfPlayAppGaming";
import GolfPlayAppFinished from "./GolfPlayApps/GolfPlayAppFinished";

type WebSocketMessage = components["schemas"]["GamePlayerMessageS2C"];

type Game = components["schemas"]["Game"];
type Problem = components["schemas"]["Problem"];

type GameState = "connecting" | "waiting" | "starting" | "gaming" | "finished";

export default function GolfPlayApp({
  game,
  sockToken,
}: {
  game: Game;
  sockToken: string;
}) {
  // const socketUrl = `wss://t.nil.ninja/iosdc-japan/2024/sock/golf/${game.game_id}/play?token=${sockToken}`;
  const socketUrl =
    process.env.NODE_ENV === "development"
      ? `ws://localhost:8002/sock/golf/${game.game_id}/play?token=${sockToken}`
      : `ws://api-server/sock/golf/${game.game_id}/play?token=${sockToken}`;

  const { sendJsonMessage, lastJsonMessage, readyState } =
    useWebSocket<WebSocketMessage>(socketUrl, {});

  const [gameState, setGameState] = useState<GameState>("connecting");

  const [problem, setProblem] = useState<Problem | null>(null);

  const [startedAt, setStartedAt] = useState<number | null>(null);

  const [timeLeftSeconds, setTimeLeftSeconds] = useState<number | null>(null);

  useEffect(() => {
    if (gameState === "starting" && startedAt !== null) {
      const timer1 = setInterval(() => {
        setTimeLeftSeconds((prev) => {
          if (prev === null) {
            return null;
          }
          if (prev <= 1) {
            clearInterval(timer1);
            setGameState("gaming");
            return 0;
          }
          return prev - 1;
        });
      }, 1000);

      const timer2 = setInterval(() => {
        const nowSec = Math.floor(Date.now() / 1000);
        const finishedAt = startedAt + game.duration_seconds;
        if (nowSec >= finishedAt) {
          clearInterval(timer2);
          setGameState("finished");
        }
      }, 1000);

      return () => {
        clearInterval(timer1);
        clearInterval(timer2);
      };
    }
  }, [gameState, startedAt, game.duration_seconds]);

  const [currentScore, setCurrentScore] = useState<number | null>(null);

  const onCodeChange = useDebouncedCallback((code: string) => {
    console.log("player:c2s:code");
    sendJsonMessage({
      type: "player:c2s:code",
      data: { code },
    });
  }, 1000);

  if (readyState === ReadyState.UNINSTANTIATED) {
    throw new Error("WebSocket is not connected");
  }

  useEffect(() => {
    if (readyState === ReadyState.CLOSING || readyState === ReadyState.CLOSED) {
      if (gameState !== "finished") {
        setGameState("connecting");
      }
    } else if (readyState === ReadyState.CONNECTING) {
      setGameState("connecting");
    } else if (readyState === ReadyState.OPEN) {
      if (lastJsonMessage !== null) {
        console.log(lastJsonMessage.type);
        if (lastJsonMessage.type === "player:s2c:prepare") {
          const { problem } = lastJsonMessage.data;
          setProblem(problem);
          console.log("player:c2s:ready");
          sendJsonMessage({ type: "player:c2s:ready" });
        } else if (lastJsonMessage.type === "player:s2c:start") {
          if (
            gameState !== "starting" &&
            gameState !== "gaming" &&
            gameState !== "finished"
          ) {
            const { start_at } = lastJsonMessage.data;
            setStartedAt(start_at);
            const nowSec = Math.floor(Date.now() / 1000);
            setTimeLeftSeconds(start_at - nowSec);
            setGameState("starting");
          }
        } else if (lastJsonMessage.type === "player:s2c:execresult") {
          const { score } = lastJsonMessage.data;
          if (
            score !== null &&
            (currentScore === null || score < currentScore)
          ) {
            setCurrentScore(score);
          }
        }
      } else {
        setGameState("waiting");
        console.log("player:c2s:entry");
        sendJsonMessage({ type: "player:c2s:entry" });
      }
    }
  }, [sendJsonMessage, lastJsonMessage, readyState, gameState, currentScore]);

  if (gameState === "connecting") {
    return <GolfPlayAppConnecting />;
  } else if (gameState === "waiting") {
    return <GolfPlayAppWaiting />;
  } else if (gameState === "starting") {
    return <GolfPlayAppStarting timeLeft={timeLeftSeconds!} />;
  } else if (gameState === "gaming") {
    return (
      <GolfPlayAppGaming
        problem={problem!.description}
        onCodeChange={onCodeChange}
        currentScore={currentScore}
      />
    );
  } else if (gameState === "finished") {
    return <GolfPlayAppFinished />;
  } else {
    return null;
  }
}
