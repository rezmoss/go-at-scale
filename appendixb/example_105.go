// Example 105
// Pitfall 1: Not checking for closed channels
func readChannel(ch chan int) int {
    return <-ch  // Panics if channel is closed!
}

// Solution: Always check for closed channels
func readChannelSafely(ch chan int) (int, bool) {
    val, ok := <-ch
    return val, ok
}

// Pitfall 2: Closing channels from the receiver
func receiverClosesPitfall(ch chan int) {
    for val := range ch {
        if isDone(val) {
            close(ch)  // Don't do this!
            return
        }
    }
}

// Solution: Let sender close the channel
func senderClosesSolution(ch chan int, done chan struct{}) {
    for val := range ch {
        if isDone(val) {
            done <- struct{}{}
            return
        }
    }
}

// Pitfall 3: Not handling channel direction properly
func confusingChannels(ch chan int) {
    // Is this a sender or receiver?
}

// Solution: Use channel direction
func sendOnly(ch chan<- int) {
    ch <- 42  // Can only send
}

func receiveOnly(ch <-chan int) {
    val := <-ch  // Can only receive
}