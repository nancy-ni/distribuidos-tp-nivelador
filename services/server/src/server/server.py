import socket
import logger
import safe_socket
from lottery import Lottery
from protocol.messages import message_codes
from protocol.messages.winner import Winner
from protocol.messages.finish import Finish
from protocol.messages.bet import BetWrapper
from protocol.communication import communication
from protocol.communication.packet import Packet


class Server:
    def __init__(self, server_host: str, server_port: int) -> None:
        self.server_host = server_host
        self.server_port = server_port

    def _handle_client(self, client_socket):
        action = "handle-client"
        message_amount = 0
        try:
            logger.info(action, logger.LogResult.in_progress)
            client_agency_id, lottery, message_amount = self.receive_bets(client_socket)

            if client_agency_id is not None:
                self.send_winners(client_socket, client_agency_id, lottery)

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

    def receive_bets(self, client_socket):
        client_agency_id = None
        lottery = None

        action = "handle-client"
        message_amount = 0
        while True:
            packet, err = communication.receive_packet(client_socket)
            if err:
                logger.error("recv-packet", logger.LogResult.fail, "messages-amount", message_amount)
                continue
            if packet.message_code == message_codes.ASK_WINNERS_CODE:
                break
            if packet.message_code != message_codes.BET_CODE:
                logger.error(action, logger.LogResult.fail, "messages-amount", message_amount)
                continue

            if client_agency_id is None:
                client_agency_id = packet.message.bet.agency_id
                lottery = Lottery(f"received_bets_{client_agency_id}.csv")
            lottery.store_bets([packet.message.bet])
            message_amount += 1

        return client_agency_id, lottery, message_amount



    def send_winners(self, client_socket, client_agency_id, lottery):
        for bet in lottery.load_bets():
            if lottery.has_won(bet):
                bet_wrapper = BetWrapper(bet)
                winner_message = Winner(bet_wrapper)
                packet = Packet(message_codes.WINNER_CODE, winner_message)
                communication.send_packet(client_socket, packet)
                print("WINNER ENVIADO CORRECTAMENTE")

        finish_message = Finish(client_agency_id)
        packet = Packet(message_codes.FINISH_CODE, finish_message)
        communication.send_packet(client_socket, packet)

    def run(self):
        action = "accept-connection"
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as server_socket:
            server_socket.bind((self.server_host, self.server_port))
            server_socket.listen()
            while True:
                try:
                    logger.info(action, logger.LogResult.in_progress)
                    client_socket, _ = server_socket.accept()
                except Exception as e:
                    logger.error(action, logger.LogResult.fail)
                    raise e
                logger.info(action, logger.LogResult.success)

                self._handle_client(client_socket)
