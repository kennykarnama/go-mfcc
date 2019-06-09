import pika
import matplotlib.pyplot as plt

credentials = pika.PlainCredentials('guest', 'guest')
connection = pika.BlockingConnection(
    pika.ConnectionParameters('localhost',5672,'/',credentials))

channel = connection.channel()

channel.queue_declare(queue='plot-me')

images = 1

def restoreData(delimitedString):
    signals = delimitedString.split(",")
    signals = [float(x) for x in signals]
    return signals

def plotSignal(signal):
    global images
    plt.figure(1)
    plt.title('Signal Wave...')
    plt.plot(signal)
    plt.savefig("signal_"+str(images)+".png")
    images = images+1

def callback(ch, method, properties, body):
    signals = restoreData(str(body, 'utf-8'))
    print("Plotting signals...")
    plotSignal(signals)


channel.basic_consume(
    queue='plot-me', on_message_callback=callback, auto_ack=True)

print(' [*] Waiting for messages. To exit press CTRL+C')
channel.start_consuming() 