#Channel

## Overview
- Channels are for orchestration between goroutines, it allows us to have 2 goroutines participate in some sort of workflow and give us the ability to orchestrate in a predictable way.
- Channel is a way of signaling event to another goroutine rather than a queue implementation.
- Three are two type of channels in golang: **buffered** and **unbuffered**.
1) Unbuffered Channel: Blocking the sender or receiver until one of directions participates in sending or receiving data
2) Buffer Channel: Don't block the sender as long as the channel's space is still available to buffer for sending. In the case, the buffer is full, the sender must be waited for( or blocked).
 At the opposite side, the receiver is blocked when **channel is empty**.
 
## Usages
- Have 5 operations work with channel: <- `chan`(receive) , chan <-(send), close, len, cap. All of them are synchronized, so no further synchronizations needed to safely perform on these operations.

    | Operation | A Nil Channel | A Closed Channel | A Not-Closed Non-Nil Channel |
    | --------------- | --------------- | --------------- | ------------------- |
    | Close | panic | panic | succeed to close |
    | Send Value To	 | block for ever | panic | block or succeed to send |
    | Receive Value From | block for ever | never block | block or succeed to receive |

- The element values of channel are transferred by copying. Three are 2 scenarios of copying value, at least on time form one goroutine to another, and 
two copies happen in transfer process if the transferred value ever stayed in the buffer of a channel. 
The size of channel element types must be smaller than 65536 bytes, the channel with larger-size types should not create to avoid overhead of copying.

- Multiple receiver goroutines can listen on the same channel, so each of them receives an *Closed* event in the case channel was closed
    ```go
    done := make(chan struct{})
    _, ok := <- done// checking on closed channel
    if !ok {
        //some code to handle
    }
    ```
- Using *range* idiom on *buffered / unbuffered channel* to receive data *until the channel is closed* , if there are multiple goroutines on the same channel, 
each of them will be received an arbitrary data until the channel is closed.
    ```go
    // the loop is terminated only when channel was closed.
    for d := range dataCh {
        // code to handle received data
    }
    ``` 
- Being a simple one to block goroutine forever through `select{}`

- Try-Send or Try-Receive with `default` case
    ```go
      ch = make(chan struct{}{}, 1)
      //try send
      for i := 0; i< 2; i++ {
          select {
              case ch<- struct{}{}:
              default:  
          }	
      }
    
      
    
      //try receive
      for i := 0; i< 2; i++ {
          select {
              case va :=<-ch:
              default:  
          }
      }
    ```
    
- Check if the channel is closed without blocking the current goroutine
    ```go
        func IsClosed(ch <- chan struct{}) {
            select {
             case <-ch: 
                return true
             default:  
            }
            return false
        }
    ```
#### Patterns
1. Double Signal: signal an event and wait for acknowledgement.
    ```go
        ch := make(chan string)
    
        go func() {
            fmt.Println(<-ch)
            ch <- "ok done"
        }()
    
        // It blocks on the receive. This Goroutine can no longer move on until we receive a signal.
        ch <- "do this"
        fmt.Println(<-ch)
    ```

1. Sending or receiving with timeout.
    > Using a buffered channel with capacity 1 to solve this issue.

    _Sending_
    ```go
        //Due to the timeout that break the waiting on `select` idiom  
        //, so it will make a deadlock if we use unbuffered channel.
        ch := make(chan struct {}, 1)
        go func() {
            time.Sleep(12 *time.Second)
            ch <- struct{}{}
        }()
        select {
            case <- ch:
                fmt.Println("received")
            case <- time.After(10 * time.Second):
                fmt.Println("timeout")	
        }
    ``` 
    _Receiving_

    ```go
        ch := make(chan struct {}, 1)
        go func() {
            time.Sleep(12 *time.Second)
            <- ch
        }()
        //Due to the timeout that break the waiting on `select` idiom  
        //, so it will make a deadlock if we use unbuffered channel.
        select {
            case ch <- struct{}{}:
                fmt.Println("received")
            case <- time.After(10 * time.Second):
                fmt.Println("timeout")	
        }
        
    ```

1. Select and Drop:
    > Walk away from a channel operation if it will immediately block

    ```go
        ch := make(chan struct {}, 1)
        go func() {
            time.Sleep(12 *time.Second)
            ch <- struct{}{}
        }()
        
        select {
        case <- ch:
            fmt.Println("received")
        case <- time.After(15 *time.Second):
            fmt.Println("timeout")
        default:
            //if we don't receive data immediately from channel
            fmt.Println("default process")
        }
        close(ch)
        
    ``` 

1. First Response Win:
    > Multiple goroutines handling the same request, but the first fastest response wins.
    >> _Solution 1:_ Using a buffered channel whose capacity is equal to number of goroutines.

    >> _Solution 2:_ Using a buffered channel whose capacity is one and select mechanism.

1. Request-Response
    > The parameters and results can be buffered channel so that the response sides won't be waited for the request sides 
    to take out the transferred value.

    > In some circumstances, using timout mechanism when the response needs longer time than the expected one to arrive.

    > A request not be guaranteed to be response back a valid value. So for all kind of reasons, an error should be returned instead.

    > A sequence of values should be returned from response sides, this is kind of data flow mechanism implementation.

1. Notifications:

    >_1 to 1_
    >> Using unbuffered or buffered channel whose capacity is one for notification between 2 goroutines

    > _1 to N or N to 1_
    >>A common usage *N To 1* is *sync.WaitGroup* and *1 to N* is unbuffered channel with *close* operation

1. Timer:

    >_Scheduled Notification_
    >> Using function *After* from package `time`

1. Lock
    > _Use Channels as Mutex Locks_
    >> One capacity buffered channel can be used as one time binary semaphore

    _Sending_
    ```go
        locker := make(chan struct{}{}, 1)
        go func() {
            locker <- struct{} {}//lock
            //handle code
            <-locker//unlock
        }()
    ```
    _Receiving_
    ```go
        locker := make(chan struct{}{}, 1)
        locker <- struct{}{}
        go func() {
            <- locker //lock
            // code
            locker <- struct{}{}//unlock
        }()
    ```

    > _Use channel as counting semaphore_
    >> Counting semaphores can be viewed as multi-owner locks. If the capacity of channel is `N`, 
    then it can be viewed as a lock which can have most N owners at any time.
    Counting semaphores are often used to enforce a maximum number of concurrent requests.

    _Sending_
    ```go
        type Seat struct {}
        type Bar chan Seat
    
        bar24x7 := make(Bar, 10)
        for i := 0; i< cap(bar24x7); i++ {
            counts <- Seat{}
        }
        for cust := 0; ; cust ++ {
            // goroutines are greater than channel capacity, 
            // they will be blocked due to reaching to capacity
            // This is inefficient 
            go func() {
                seat := <- bar247 
                //code
                bar247 <- seat
            }()
        }
        
        
    ```

    _Receiving_
    
    ```go
        type Seat struct {}
        type Bar chan Seat
    
        bar24x7 := make(Bar, 10)
        for i := 0; i< cap(bar24x7); i++ {
            counts <- Seat{}
        }
        //at most 10 goroutines coexisting, this is optimized version in comparision with above solution 
        for cust := 0; ; cust ++ {
            seat := <- bar247 
            go func(s * Seat) {
                //code
                bar247 <- *seat
            }(&seat)
        }
        
        
    ```

1. Dialogue(Ping-Pong)
    > Using an unbuffered channel to implement this pattern

    ```go
    func Play(playName string, ball chan) {
        for {
            ballNum := <- ball
            ballNum ++
            if ballNum < 100 {
                return
            }
            ball <- ballNum
        }
    }
    
    ballCh := make(chan uint64)
    
    go Play("A", ballCh)
    Play("B", ballCh)
    ```

 

 
#### Stop gracefully

## References
- https://github.com/hoanhan101/ultimate-go/blob/master/go/concurrency/channel_1.go
- https://go101.org/article/channel.html
- https://go101.org/article/channel-use-cases.html
- https://go101.org/article/channel-closing.html