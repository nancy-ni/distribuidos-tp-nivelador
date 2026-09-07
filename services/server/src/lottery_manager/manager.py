import queue
import threading

# Entidad encargada de gestionar las apuestas y el sorteo mismo.
# Los threads handlers de cada agencia se comunican para reportar las apuestas recibidas, y
# luego cuando el LotteryManager hace el sorteo, les responde los respectivos ganadores.
class LotteryManager:
    def __init__(self, lottery, min_quorum):
        self.lottery = lottery
        self.min_quorum = min_quorum
        self.ready_queue = queue.Queue()
        self.shutdown_event = threading.Event()
        self.thread_handler = threading.Thread(target=self._run)
        self.lock = threading.Lock()


    def start(self):
        self.thread_handler.start()


    def stop(self):
        self.shutdown_event.set()
        self.ready_queue.put(None)
        self.thread_handler.join()


    def store_bets(self, bets):
        with self.lock:
            self.lottery.store_bets(bets)


    # Utilizada cuando un thread handler de agencia quiere indicar que esta lista para los 
    # resultados. Devuelve un nuevo canal a traves del cual recibira los ganadores.
    def report_ready(self, agency_id):
        response_queue = queue.Queue()
        self.ready_queue.put((agency_id, response_queue))
        return response_queue


    def _run(self):
        ready_agencies = {}

        while not self.shutdown_event.is_set():
            ready_report = self.ready_queue.get()
            if ready_report is None:
                break

            agency_id, response_queue = ready_report
            ready_agencies[agency_id] = response_queue

            if len(ready_agencies) >= self.min_quorum:
                self._send_all_winners(ready_agencies)
                ready_agencies.clear()

        for _, response_queue in ready_agencies.items():
            response_queue.put(None)


    # Realiza el sorteo y les comunica a c/u de los threads handlers de agencias listas sus 
    # respectivos ganadores, si tienen.
    def _send_all_winners(self, ready_agencies):
        with self.lock:
            for bet in self.lottery.load_bets():
                agency_id = bet.agency_id
                if self.lottery.has_won(bet) and agency_id in ready_agencies:
                    response_queue = ready_agencies[agency_id]
                    response_queue.put(bet)
            
        for _, response_queue in ready_agencies.items():
            response_queue.put(None)
