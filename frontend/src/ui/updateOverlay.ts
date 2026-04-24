import { reactive } from 'vue';
import { getWsUrl } from '../api';
import { updates } from './updates';

const stripAnsi = (text: string) => text.replace(/\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])/g, '');

const state = reactive({
  visible: false,
  output: '',
  follow: true,
  waitingForReload: false,
  progressFailed: false,
});

let socket: WebSocket | null = null;

const append = (text: string) => {
  state.output += stripAnsi(text);
};

const disconnect = () => {
  if (socket) {
    socket.close();
    socket = null;
  }
};

const connect = (messages: { socketError: string; socketClosed: string }) => {
  disconnect();
  socket = new WebSocket(getWsUrl('/app-updates'));
  socket.onmessage = (event) => append(String(event.data));
  socket.onerror = () => append(`\n[error] ${messages.socketError}\n`);
  socket.onclose = () => {
    socket = null;
    append(`\n[closed] ${messages.socketClosed}\n`);
    void updates.syncStatus().catch(() => undefined);
  };
};

const open = (messages: { socketError: string; socketClosed: string }, opts?: { resetOutput?: boolean }) => {
  state.visible = true;
  if (!socket) {
    if (opts?.resetOutput) {
      state.output = '';
    }
    state.follow = true;
    connect(messages);
  }
};

const close = () => {
  state.visible = false;
  disconnect();
};

const beginUpdate = (messages: { socketError: string; socketClosed: string; starting: string }) => {
  open(messages, { resetOutput: true });
  state.progressFailed = false;
  state.waitingForReload = false;
  append(`[info] ${messages.starting}\n`);
};

const waitForReload = (messages: { reloadTimeout: string }) => {
  state.waitingForReload = true;
  void updates.waitForAppReload().then((reloaded) => {
    if (!reloaded) {
      state.waitingForReload = false;
      state.progressFailed = true;
      append(`\n[warn] ${messages.reloadTimeout}\n`);
    }
  });
};

const markFailed = (message: string) => {
  state.waitingForReload = false;
  state.progressFailed = true;
  append(`[error] ${message}\n`);
};

export const updateOverlay = {
  state,
  append,
  connect,
  disconnect,
  open,
  close,
  beginUpdate,
  waitForReload,
  markFailed,
};
