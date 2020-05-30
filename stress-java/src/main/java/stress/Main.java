package stress;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufAllocator;
import io.netty.buffer.PooledByteBufAllocator;

import java.util.LinkedList;

public class Main {
    public static void main(String[] args) {
        LinkedList list = new LinkedList();

        final int chunkSizeM = 20;
        final int maxTotalSizeM = 512;
        final long waitSec = 1;

        for(int currentTotalSizeM = 0; currentTotalSizeM + chunkSizeM < maxTotalSizeM;) {
            ByteBufAllocator alloc = PooledByteBufAllocator.DEFAULT;

            ByteBuf buf = alloc.directBuffer(chunkSizeM*1024*1024 );
            buf.setZero(0, buf.capacity());
            currentTotalSizeM += chunkSizeM;
            System.out.println("currentTotalSizeM = " + currentTotalSizeM);

//        buf.release(); // The direct buffer is returned to the pool.
            list.add(buf);
            try {
                Thread.sleep( waitSec * 1000 );
            } catch (InterruptedException e) {
                e.printStackTrace();
            }

        }

        System.out.println("Stress java exit.");

    }
}
