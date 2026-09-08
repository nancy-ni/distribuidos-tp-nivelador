# 1. Introducción

En el presente trabajo se implementó un protocolo de comunicación entre cliente y servidor para modelar un sistema distribuido simple inspirado en la central de Lotería Nacional. En dicho sistema, distintas agencias de lotería cargan participantes en un servidor central, que realizará el sorteo e informará los participantes que han ganado.

A continuación, se detalla brevemente el protocolo de comunicación implementado y los mecanismos para sincronizar la ejecución concurrente.


# 2. Protocolo de Comunicación

## 2.1 Flujo General

El flujo central del sistema empieza cuando, tras la conexión del cliente al servidor, el cliente comienza a leer desde el archivo de input especificado y, a medida que va leyendo, extrae los datos de las distintas apuestas para mandarlas en batches. Por cada batch enviado, se espera una confirmación del servidor antes de seguir. Cuando se terminan de enviar todas las apuestas, el cliente envía un mensaje especial para peticionar obtener los ganadores del sorteo. Luego, el cliente se queda esperando para recibir tales ganadores y, al recibir uno, escribe sus datos en el archivo de output especificado. Cabe destacar que el servidor sólo realiza el sorteo cuando llega a un mínimo quórum con las demás agencias. Finalmente, el intercambio termina cuando el servidor envía un mensaje especial de fin, tras el cual se cierra su conexión limpiamente.

En resumen, las etapas ó pasos principales en la comunicación entre cliente (agencias) y servidor son, desde la perspectiva del cliente/servidor respectivamente:
* Envío / Recepción de **apuestas** en batches.
* Envío / Recepción de **petición** para obtener los ganadores.
* Recepción / Envío de **ganadores**.
* Recepción / Envío de mensaje de indicación de **fin**.

Por ejemplo, para el caso con un cliente y un quórum mínimo de 1, el diagrama de secuencia simplificado sería:
<p align="center">
  <img width="491" height="636" alt="Captura de pantalla 2026-09-08 010740" src="https://github.com/user-attachments/assets/3a0eed7b-6d8a-4986-9d14-8f5fe79369a4" />
</p>


## 2.2 Comunicación

El cliente y el servidor se comunican a través de la transmisión e intercambio de paquetes. Cada paquete cuenta con un campo `message_code` y `message`.
* El campo `message_code` se utiliza para distinguir los distintos tipos de mensajes internos del protocolo que se envían entre sí. Cada mensaje tiene su propio código. Ej: BATCH_CODE=0x02
* El campo `message` es ocupado por el mensaje en sí. Este mensaje fue implementado como una interfaz posible de serializar y deserializar. Cada mensaje intercambiado entre cliente y servidor cumple con esta interfaz.

Para la transmisión de paquetes, como cada uno posee un largo variable dependiendo del mensaje interno que lleva, se utiliza un modelo de bloques dinámicos donde se agrega la longitud del paquete como prefijo al mismo, en particular, se reservan 2 bytes para especificar tal longitud. Esto permite al otro extremo saber la cantidad exacta de bytes a leer cuando espera recibir un paquete.

Gráficamente, lo que se envía / recibe por el stream de bytes como paquete es:

<p align="center">
  <img width="260" height="185" align="center" alt="Captura de pantalla 2026-09-07 040339" src="https://github.com/user-attachments/assets/ce485029-6288-4de2-9b5d-34f2be7484bc" />
</p>

_Nota_: En el gráfico no se tiene en cuenta el tamaño en bytes de los campos.


## 2.3 Mensajes Intercambiados

A continuación se presentan los principales mensajes intercambiados durante la comunicación.

### Batch

* _Emisor_: Cliente / Agencia.
* _Receptor_: Servidor de lotería.
* _Propósito_: Comunicación de apuestas al servidor. Los batches permiten acortar tiempos de transmisión y procesamiento en comparación con mandar sólo una apuesta a la vez.

```go
type Batch struct {
	Bets []Bet
}
```
```go
type Bet struct {
	agencyId  uint32
	firstName string
	lastName  string
	dni       uint32
	birthday  string
	betNumber uint16
}
```

* Para serializar una apuesta (`Bet`), se utilizan un modelo de bloques fijos para los datos que poseen longitudes constantes (`agency_id`, `dni`, `bet_number`). Para los datos con largo variable (`firstname`, `lastname`), se agrega como prefijo la longitud del campo previo al dato en sí, en este caso, se utiliza 1 byte para describir la longitud del nombre y apellido. Finalmente, para el caso especial de la fecha de nacimiento, se decidió serializarlo como string con longitud fija, siguiendo el formato constante de YYYY-MM-DD (10 bytes en total) por simplicidad.

* Para serializar un batch (`Batch`), se incluye como prefijo la cantidad de apuestas contenidas dentro del batch, para tener en consideración el caso donde no se llegue a completar un batch entre todos los apostantes. Luego, por cada apuesta, se agrega como prefijo la longitud de la apuesta total, para poder demarcar el inicio y fin de cada apuesta en el stream de bytes.

Gráficamente, lo que se envía / recibe por el stream de bytes dentro de un paquete, al hablar de un `Batch` y `Bet` respectivamente, será:
<p align="center">
  <img width="670" height="397" alt="Captura de pantalla 2026-09-07 033119" src="https://github.com/user-attachments/assets/392e570d-6c82-4062-8582-0924c45b93b6" />
</p>

_Nota_: En el gráfico no se tiene en cuenta el tamaño en bytes de los campos.


### Ack

* _Emisor_: Servidor de lotería.
* _Receptor_: Cliente / Agencia.
* _Propósito_: Confirmar la correcta recepción del batch completo. También sirven para sincronizar las tasas de envío y recepción de paquetes entre cliente y servidor, sumando soporte para el control de congestión best effort que ya posee TCP.

```go
type Ack struct {
	agencyId uint32
}
```


### AskWinners

* _Emisor_: Cliente / Agencia.
* _Receptor_: Servidor de lotería.
* _Propósito_: Peticionar obtener la lista de ganadores para el respectivo cliente / agencia. Cuando el servidor recibe este mensaje, deja de escuchar por nuevos mensajes de batches, y suma al cliente / agencia como uno de los que aportan al quórum para realizar el sorteo.

```go
type AskWinners struct {
	agencyId uint32
}
```


### Finish

* _Emisor_: Servidor de lotería.
* _Receptor_: Cliente / Agencia.
* _Propósito_: Indica el fin del intercambio entre el cliente y el servidor. Cuando el servidor termina de enviar este mensaje, procede a cerrar limpiamente la conexión con el cliente. Análogamente, cuando el cliente recibe este mensaje, deja de escuchar por nuevos ganadores y cierra la conexión.

```go
type Finish struct {
	agencyId uint32
}
```

Cabe destacar que en este protocolo el servidor, para comunicar las apuestas ganadoras a las agencias, reutiliza el mensaje de `Bet` por simplicidad, sin crear un nuevo tipo de mensaje específico para comunicar los ganadores (aunque sí tiene un `message_code` distinguible).



# 3. Concurrencia y Sincronización

## 3.1 Modelo Elegido

Para que el servidor pueda tratar con múltiples clientes de manera concurrente, se eligió un modelo de **multithreading**. 

Considerando que el servidor está escrito en Python, es notable destacar que el lenguaje posee el _Global Interpreter Lock (GIL)_: un mutex que protege el acceso a los objetos internos de Python, permitiendo que sólo un hilo se ejecute a la vez. Este comportamiento suele ser perjudicial para programas multithreading que sean CPU intensive, pues incluso si se utilizan hilos para dividir las tareas, sólo uno puede avanzar productivamente en ellas, quitando el beneficio principal de tener threads.

Sin embargo, en este caso, el GIL no resulta un problema para lograr la concurrencia: en este sistema, las operaciones son más orientadas a I/O que a tareas pesadas. El cliente y el servidor pasan la mayor parte del tiempo de ejecución enviando y recibiendo paquetes, lo cual permite al sistema operativo intercalar otras tareas durante los tiempos de espera. De esta manera, a pesar de las restricciones del GIL, el servidor es capaz de manejar las múltiples conexiones con distintos clientes de manera concurrente.

Teniendo esto en cuenta, a continuación se detallan las herramientas y mecanismos utilizados en la implementación para lograr la sincronización entre los distintos hilos.

## 3.2 Implementación Concurrente

Para manejar múltiples conexiones con los distintos clientes, el servidor lanza un hilo por cada uno de ellos, para manejarlos de manera independiente. Para luego poder cerrar los sockets cuando finaliza el intercambio, el servidor guarda las conexiones en un arreglo protegido bajo un Lock.

Cuando un hilo recibe un batch de apuestas, éste debe guardarlo a través del método `store_bets` de la clase `Lottery` provista, sin embargo, los métodos de esta clase no son thread-safe. Para resolver este problema, se implementó una entidad central `LotteryManager` la cual se ejecuta en un hilo aparte por su cuenta, éste tiene las siguientes funcionalidades:
* _Controla la clase Lottery bajo un lock_: `LotteryManager` cuenta con un método que permite a los threads handlers tomar el lock para la clase `Lottery` y así guardar las apuestas, sin riesgo de race conditions.
* _Recibe la notificación cuando el cliente está listo_: Cuando el thread handler del cliente recibe el mensaje para obtener los ganadores, el thread se lo comunica al `LotteryManager` a través de un canal (queue en Python) que comparten todos los threads. El `LotteryManager` lleva cuenta de cuántos clientes hacen falta para llegar al quórum, y cuando recibe la notificación, devuelve otro canal distinto al thread por donde comunicará eventualmente los ganadores si tuviese.
* _Realiza el sorteo_: Cuando el servidor recibe el mínimo de notificaciones para llegar al quórum, `LotteryManager` se encarga de iterar entre las apuestas cargadas e identificar los ganadores pertenecientes a los clientes listos.
* _Comunica los ganadores_: Cuando el `LotteryManager` identifica un ganador, se encarga de comunicarlo al thread handler del cliente al cual pertenece el ganador, mediante un canal (queue) específico entre el manager y el hilo, mencionado previamente.

```python
class LotteryManager:
    def __init__(self, lottery, min_quorum):
        self.lottery = lottery
        self.min_quorum = min_quorum
        self.ready_queue = queue.Queue()
        self.shutdown_event = threading.Event()
        self.thread_handler = threading.Thread(target=self._run)
        self.lock = threading.Lock()

    def store_bets(self, bets):
        with self.lock:
            self.lottery.store_bets(bets)
```

En resumen, dentro del servidor existe: un hilo dedicado a **aceptar conexiones** con los clientes, un hilo para el `LotteryManager` encargado de la **lógica central** para manejar las apuestas, y un hilo por **cada conexión** con un cliente.
<p align="center">
  <img width="802" height="97" alt="Captura de pantalla 2026-09-08 001104" src="https://github.com/user-attachments/assets/6d704b55-7b85-434b-93ab-f42e98a67aef" />
</p>

Luego, en particular, las herramientas que utiliza el `LotteryManager` para lograr la sincronización entre los distintos hilos son:
* **Lock** para proteger la sección crítica que implica utilizar la clase `Lottery` para guardar las apuestas.
* **Canales** (Queues en Python):
  * **Canal único** para que los distintos threads handlers de clientes puedan notificar al `LotteryManager` que están listos para hacer el sorteo (luego de guardar todas las apuestas que el cliente transmitió). Debido a la naturaleza de los canales, tener un único canal compartido entre distintos hilos no genera race conditions.
  * **Un canal por thread**, que crea el propio `LotteryManager` y se lo entrega al thread cuando éste envía la notificación. Cuando eventualmente se haga el sorteo, se comunican los ganadores de cada agencia por acá.
<p align="center">
  <img width="787" height="335" alt="Captura de pantalla 2026-09-07 204932" src="https://github.com/user-attachments/assets/cd3c83d9-818c-4af3-9eef-296228d56a19" />
</p>


## 3.3 Graceful Shutdown

### 3.3.1 Servidor

Cuando el intercambio entre cliente y servidor termina normalmente sin interrupciones, los threads handlers del servidor acceden concurrentemente, tomando el lock del servidor, a la lista de conexiones para cerrar el socket correspondiente.

Para tratar el caso donde el servidor recibe una señal _SIGTERM_, se utilizó la librería signal para capturar dicha señal, y se hace uso de eventos de la librería threading para activar una bandera de cierre y así notificar al resto del sistema. En particular, las acciones que siguen a la activación de la bandera son:
* Cierre del socket de escucha del servidor.
* Cierre del `LotteryManager`:
  * Se envía un mensaje vacío (ej. None) a través del canal `ready_queue`, permitiendo salir de la espera por notificaciones.
  * Se envían mensajes vacíos (ej. None) a través de los canales `response_queue`, permitiendo a los threads handlers salir de la espera por los ganadores.
  * Se espera a joinear el thread del `LotteryManager`.
* Cierre de los sockets de las conexiones con los clientes.
* Se espera a joinear todos los threads handlers de clientes.

### 3.3.2 Cliente

Similarmente del lado del cliente, existe una bandera para indicar el cierre de la ejecución. El acceso a tal bandera está protegida bajo un Lock. 

Para manejar la señal de _SIGTERM_, se utiliza la librería os/signal para crear un canal que intercepte la interrupción. Además del hilo principal del cliente, se lanza una go routine para escuchar sobre tal canal y, en caso de recibir la señal, se activa la bandera de cierre antes de proceder a cerrar el socket de la conexión.

Para evitar que la go routine persista incluso cuando la conexión cierra limpiamente sin interrupciones, existe otro canal simple tal que emite un valor cuando termina el intercambio, permitiendo que la go routine retorne.





