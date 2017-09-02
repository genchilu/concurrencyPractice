package url.genchi.lmax;

import com.lmax.disruptor.EventHandler;

/**
 * Created by mac on 2017/4/21.
 */
public class LongEventHandler implements EventHandler<LongEvent>
{
    public void onEvent(LongEvent event, long sequence, boolean endOfBatch)
    {
        System.out.println("Event: " + event);
    }
}