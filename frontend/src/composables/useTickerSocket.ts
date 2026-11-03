import { ref, shallowRef, onMounted, onUnmounted } from 'vue'
import { Tick } from '../gen/ticker'

export function useTickerSocket(url = 'ws://localhost:8080/ws') {
  const isConnected = ref(false)
  // shallowRef dipakai agar Vue tidak membuat proxy mendalam untuk tiap field setiap 500ms
  const ticks = shallowRef<Record<string, Tick>>({})
  let ws: WebSocket | null = null

  const connect = () => {
    ws = new WebSocket(url)
    // WAJIB: ubah mode terimaan dari 'blob' menjadi 'arraybuffer' untuk Protobuf
    ws.binaryType = 'arraybuffer'

    ws.onopen = () => {
      isConnected.value = true
    }

    ws.onmessage = (event: MessageEvent<ArrayBuffer>) => {
      try {
        const bytes = new Uint8Array(event.data)
        const decoded = Tick.decode(bytes)

        ticks.value = {
          ...ticks.value,
          [decoded.symbol]: decoded,
        }
      } catch (err) {
        console.error('Gagal decode binary Protobuf:', err)
      }
    }

    ws.onclose = () => {
      isConnected.value = false
    }

    ws.onerror = (err) => {
      console.error('WebSocket error:', err)
      ws?.close()
    }
  }

  onMounted(() => connect())
  onUnmounted(() => ws?.close())

  return { isConnected, ticks }
}
