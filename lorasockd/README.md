# lorasockd

Lorasockd is a daemon to provide an API to a SX1278-based SPI Lora Module (AI Thinker RA-02) via Unix Socket. While it does create more to maintain, this application enables developers to use Lora as any other Linux utility, via any FS/Unix Socket capable language/tool. Additionally, Providing this API over Unix Socket makes it very easy to fake/mock on non-embedded development devices, Greatly improving developer workflow.
