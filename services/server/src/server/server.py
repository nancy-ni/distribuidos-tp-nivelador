import socket
import logger
import os
import threading
from lottery import Lottery
from protocol.messages import message_codes
from protocol.messages.winner import Winner
from protocol.messages.finish import Finish
from protocol.messages.bet import BetWrapper
from protocol.communication import communication
from protocol.communication.packet import Packet
from lottery_manager.manager import LotteryManager


class Server:
    def __init__(self, server_host: str, server_port: int) -> None:
        self.server_host = server_host
        self.server_port = server_port
        self.shutdown_event = threading.Event()

    def _handle_client(self, client_socket, lottery_manager):
        action = "handle-client"
        message_amount = 0
        try:
            logger.info(action, logger.LogResult.in_progress)
            client_agency_id, message_amount = self.receive_bets(client_socket, lottery_manager)

            if client_agency_id is not None:
                self.send_winners(client_socket, client_agency_id, lottery_manager)

            logger.info(
                action,
                logger.LogResult.success,
                "messages-amount",
                message_amount,
            )

        except Exception as e:
            logger.error(
                action, logger.LogResult.fail, "messages-amount", message_amount
            )
            raise e

    def receive_bets(self, client_socket, lottery_manager):
        client_agency_id = None

        action = "handle-client"
        message_amount = 0
        while True:
            packet, err = communication.receive_packet(client_socket)
            if err:
                logger.error("recv-packet", logger.LogResult.fail, "messages-amount", message_amount)
                continue
            if packet.message_code == message_codes.ASK_WINNERS_CODE:
                break
            if packet.message_code != message_codes.BATCH_CODE or len(packet.message.bets) == 0:
                logger.error(action, logger.LogResult.fail, "messages-amount", message_amount)
                continue

            if client_agency_id is None:
                client_agency_id = packet.message.bets[0].bet.agency_id

            lottery_manager.store_bets(packet.message.get_bets())
            message_amount += 1

        return client_agency_id, message_amount

    def send_winners(self, client_socket, client_agency_id, lottery_manager):
        response_queue = lottery_manager.report_ready(client_agency_id)

        while not self.shutdown_event.is_set():
            winner_bet = response_queue.get()
            if winner_bet is None:
                break

            winner_bet_wrapper = BetWrapper(winner_bet)
            winner_message = Winner(winner_bet_wrapper)
            packet = Packet(message_codes.WINNER_CODE, winner_message)
            communication.send_packet(client_socket, packet)
            print("WINNER ENVIADO CORRECTAMENTE")

        finish_message = Finish(client_agency_id)
        packet = Packet(message_codes.FINISH_CODE, finish_message)
        communication.send_packet(client_socket, packet)

    def run(self):
        action = "accept-connection"
        
        lottery = Lottery("received_bets.csv")
        min_quorum = int(os.getenv("AGENCY_QUORUM_MIN"))
        lottery_manager = LotteryManager(lottery, min_quorum)
        lottery_manager.start()

        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as server_socket:
            server_socket.bind((self.server_host, self.server_port))
            server_socket.listen()

            handlers = []
            while True:
                try:
                    logger.info(action, logger.LogResult.in_progress)
                    client_socket, _ = server_socket.accept()
                except Exception as e:
                    logger.error(action, logger.LogResult.fail)
                    raise e
                logger.info(action, logger.LogResult.success)

                thread = threading.Thread(target=self._handle_client, args=(client_socket, lottery_manager))
                handlers.append(thread)

                thread.start()
