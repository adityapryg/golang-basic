# 21 – Concurrency di Go

Go dirancang dengan dukungan concurrency tingkat pertama. Bagian ini merangkum konsep inti, primitive, pola umum, dan praktik terbaik untuk menulis program concurrent yang aman dan efisien.

## Konsep Dasar

- Concurrency vs Parallelism: concurrency adalah mengelola banyak pekerjaan yang berjalan saling tumpang-tindih; parallelism menjalankan banyak pekerjaan secara simultan pada beberapa CPU core.
- Goroutine: unit eksekusi ringan yang dikelola oleh runtime Go; pembuatan goroutine sangat murah.
- Channel: jalur komunikasi antar goroutine yang aman; mendukung pengiriman dan penerimaan nilai.
- Select: multiplexer untuk membaca dari beberapa channel sekaligus.

## Goroutine

- Membuat goroutine menggunakan kata kunci `go`.
- Scheduler Go memetakan ribuan goroutine ke thread OS secara efisien.
- `runtime.GOMAXPROCS` mengatur jumlah CPU logical yang digunakan untuk eksekusi goroutine.

```go
package main
import "fmt"
func main(){
    ch := make(chan int)
    go func(){ 
        ch <- 42 
    }()
    v := <-ch
    fmt.Println(v)
}
```

Flow:
- Buat `ch` sebagai channel integer.
- Jalankan goroutine yang mengirim `42` ke `ch`.
- Goroutine utama menerima dari `ch` dan mencetak nilai.

## Channel

- Unbuffered: operasi kirim/terima bersifat sinkron; pengirim menunggu penerima dan sebaliknya.
- Buffered: memiliki kapasitas; pengirim tidak selalu menunggu selama buffer belum penuh.
- Penutupan: `close(ch)` memberi sinyal tidak ada lagi nilai; penerima dapat memakai `for range`.
- Arah channel: `chan<- T` hanya kirim, `<-chan T` hanya terima.

```go
package main
import "fmt"
func main(){
    ch := make(chan int, 2)
    ch <- 1
    ch <- 2
    close(ch)
    for v := range ch { 
        fmt.Println(v)
    }
}
```

Flow:
- Buat channel buffered dengan kapasitas 2.
- Kirim dua nilai ke dalam buffer tanpa blocking.
- Tutup channel untuk memberi sinyal tidak ada nilai lagi.
- Iterasi `for range` membaca semua nilai hingga channel habis.

## Select

- Memilih case yang siap dieksekusi dari beberapa operasi channel.
- Mendukung timeout dan non-blocking dengan `default`.

```go
package main
import (
    "fmt"
    "time"
)
func main(){
    ch := make(chan int)
    select {
    case v := <-ch:
        fmt.Println(v)
    case <-time.After(100 * time.Millisecond):
        fmt.Println("timeout")
    }
}
```

Flow:
- Siapkan channel `ch` yang belum ada pengirim.
- `select` menunggu case yang siap: baca dari `ch` atau timeout via `time.After`.
- Karena tidak ada kiriman, case timeout dieksekusi dan mencetak "timeout".

## Sinkronisasi

- `sync.WaitGroup` menunggu kumpulan goroutine selesai.
- `sync.Mutex` dan `sync.RWMutex` melindungi data bersama dari race.
- `sync.Once` memastikan sebuah fungsi hanya dieksekusi sekali.
- `sync.Cond` untuk koordinasi lebih kompleks.
- `sync/atomic` menyediakan operasi atomik pada tipe numerik dan pointer.

```go
package main
import (
    "fmt"
    "sync"
)
func main(){
    var wg sync.WaitGroup
    jobs := []int{1,2,3,4}
    wg.Add(len(jobs))
    for _, j := range jobs {
        go func(x int){ 
            defer wg.Done(); 
            fmt.Println(x)
        }(j)
    }
    wg.Wait()
}
```

Flow (WaitGroup):
- Inisialisasi `WaitGroup` dan set jumlah pekerjaan.
- Luncurkan goroutine untuk setiap job, panggil `Done()` saat selesai.
- Panggil `Wait()` di goroutine utama hingga semua selesai.

```go
package main
import (
    "fmt"
    "sync"
)
func main(){
    var mu sync.Mutex
    var n int
    var wg sync.WaitGroup
    wg.Add(2)
    go func(){ 
        defer wg.Done(); 
        for i:=0;i<1000;i++{ 
            mu.Lock(); 
            n++; 
            mu.Unlock()
        }
    }()
    go func(){ 
        defer wg.Done(); 
        for i:=0;i<1000;i++{ 
            mu.Lock(); 
            n++; 
            mu.Unlock()
        } 
    }()
    wg.Wait()
    fmt.Println(n)
}
```

Flow (Mutex):
- Dua goroutine mengakses counter bersama `n`.
- Setiap increment dibungkus `mu.Lock()/mu.Unlock()` untuk mencegah race.
- Setelah kedua goroutine selesai, nilai `n` konsisten.

```go
package main
import (
    "fmt"
    "sync/atomic"
    "sync"
)
func main(){
    var n int64
    var wg sync.WaitGroup
    wg.Add(2)
    go func(){ 
        defer wg.Done(); 
        for i:=0;i<1000;i++{ 
            atomic.AddInt64(&n,1) 
        }
    }()
    go func(){ 
        defer wg.Done(); 
        for i:=0;i<1000;i++{ 
            atomic.AddInt64(&n,1)
        }
    }()
    wg.Wait()
    fmt.Println(n)
}
```

Flow (Atomic):
- Dua goroutine menambah `n` menggunakan `atomic.AddInt64`.
- Operasi atomik memastikan penambahan aman tanpa mutex.

## Context

- `context` membawa sinyal cancel, deadline, dan nilai ke dalam call chain.
- Dipakai untuk membatalkan pekerjaan, memberikan timeout, dan propagasi informasi request.

```go
package main
import (
    "context"
    "fmt"
    "time"
)
func main(){
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()
    ch := make(chan string)
    go func(){ 
        select{ 
        case <-ctx.Done(): 
            return; 
        case ch <- "ok": 
        }
    }()
    select {
    case s := <-ch:
        fmt.Println(s)
    case <-ctx.Done():
        fmt.Println("cancel")
    }
}
```

Flow:
- Buat `context.WithTimeout` dengan batas waktu tertentu.
- Jalankan goroutine yang mengirim nilai kecuali `ctx.Done()` aktif lebih dulu.
- `select` di goroutine utama memilih antara menerima nilai atau menangani cancel/timeout.

## Pola Umum

- Worker Pool: batasi jumlah goroutine yang memproses pekerjaan.
- Pipeline: rangkaian tahap pemrosesan data melalui channel.
- Fan-in/Fan-out: gabungkan banyak channel atau sebarkan pekerjaan ke banyak worker.
- Bounded Concurrency: gunakan semaphore atau buffered channel untuk membatasi tingkat concurrency.
- Rate Limiting: gunakan `time.Ticker` atau token bucket sederhana.

```go
package main
import (
    "fmt"
    "sync"
)
func main(){
    jobCh := make(chan int)
    resCh := make(chan int)
    var wg sync.WaitGroup
    for i:=0;i<3;i++{
        wg.Add(1)
        go func(){ 
            defer wg.Done(); 
            for j := range jobCh { 
                resCh <- j*j
            } 
        }()
    }
    go func(){ 
        for i:=1;i<=5;i++{ 
            jobCh <- i 
        } 
        close(jobCh) 
    }()
    go func(){ 
        wg.Wait(); 
        close(resCh)
    }()
    for r := range resCh { 
        fmt.Println(r)
    }
}
```

Flow (Worker Pool):
- `jobCh` menyalurkan tugas, `resCh` menyalurkan hasil.
- Luncurkan beberapa worker yang membaca dari `jobCh`, memproses, lalu mengirim hasil ke `resCh`.
- Producer mengirim sejumlah job dan menutup `jobCh`.
- Goroutine penutup menunggu semua worker selesai lalu menutup `resCh`.
- Goroutine utama `range` di `resCh` untuk mengonsumsi hasil.

## Pitfall & Praktik Terbaik

- Hindari race condition dengan Mutex atau operasi atomik.
- Hindari deadlock: pahami urutan penguncian dan interaksi channel.
- Tutup channel dari pengirim, bukan penerima.
- Jangan kirim ke channel yang sudah ditutup.
- Hindari goroutine leak: pastikan setiap goroutine bisa berhenti, gunakan `context` atau sinyal penutupan.
- Gunakan `go test -race` untuk mendeteksi data race.
- Batasi jumlah goroutine; gunakan worker pool atau bounded concurrency.

## Timeout, Ticker, dan Time-based Control

- `time.After` untuk timeout satu kali.
- `time.Ticker` untuk event berkala.
- Kombinasikan dengan `select` untuk kontrol waktu yang robust.

## Ringkasan

- Gunakan goroutine untuk menjalankan pekerjaan bersamaan.
- Komunikasikan antar goroutine dengan channel dan `select`.
- Lindungi data bersama dengan `sync` dan `atomic`.
- Kelola pembatalan dan batas waktu dengan `context`.
- Terapkan pola worker pool, pipeline, dan fan-in/fan-out untuk arsitektur yang bersih.