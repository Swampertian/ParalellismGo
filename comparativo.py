import requests
import time
import multiprocessing as mp
import os

def fetch_endpoint1(queue):
    start = time.time()
    url = "https://apitempo.inmet.gov.br/estacoes/T"
    try:
        response = requests.get(url)
        data = response.json()
        elapsed = time.time() - start
        print(f"Tempo de execução de fetch_endpoint1: {elapsed:.6f}s")
        queue.put(data)
    except Exception as e:
        print(f"Erro em fetch_endpoint1: {e}")
        queue.put([])

def fetch_endpoint2(queue):
    start = time.time()
    token = os.environ.get("INMET_TOKEN")
    if not token:
        print("Erro em fetch_endpoint2: INMET_TOKEN environment variable not set")
        queue.put([])
        return
    url = f"https://apitempo.inmet.gov.br/token/estacao/diaria/2022-11-01/2022-11-01/A001/{token}"
    try:
        response = requests.get(url)
        data = response.json()
        elapsed = time.time() - start
        print(f"Tempo de execução de fetch_endpoint2: {elapsed:.6f}s")
        queue.put(data)  
    except Exception as e:
        print(f"Erro em fetch_endpoint2: {e}")
        queue.put([])

if __name__ == "__main__":
    start_total = time.time()

    # Usar Queue para comunicação entre processos
    queue1 = mp.Queue()
    queue2 = mp.Queue()

    p1 = mp.Process(target=fetch_endpoint1, args=(queue1,))
    p2 = mp.Process(target=fetch_endpoint2, args=(queue2,))

    p1.start()
    p2.start()

    # Aguardar e obter resultados
    estacoes = queue1.get()
    dados = queue2.get()

    # Aguardar processos terminarem
    p1.join()
    p2.join()

    elapsed_total = time.time() - start_total
    print(f"Tempo total de execução: {elapsed_total:.6f}s")

    print("=== ESTACOES ===")
    for e in estacoes:
        print(f"Código: {e.get('CD_OSCAR', 'N/A')}")

    print("=== DADOS ===")
    for d in dados:
        print(f"Temperatura Minima : {d.get('TEMP_MIN', 'N/A')}")