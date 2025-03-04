import ws from "k6/ws";
import { sleep } from "k6";
import { replit } from "k6/x/replit";

export default function () {
  sleep(0.5); // Wait for the WS server to start

  const url = "ws://localhost:8080/repl";
  const res = ws.connect(url, null, function (socket) {
    socket.on("message", function(data) {
      try {
        let result = replit.run(data);
        socket.send(eval(data));
      } catch (error) {
        socket.send(error.toString());
      }
    });
  });
}
