package com.onepick.idem

import jakarta.servlet.http.HttpServletRequest
import jakarta.servlet.http.HttpServletResponse
import jakarta.servlet.http.HttpServletResponseWrapper
import org.slf4j.LoggerFactory
import org.springframework.stereotype.Component
import org.springframework.web.servlet.AsyncHandlerInterceptor
import java.io.ByteArrayOutputStream
import java.io.PrintWriter
import java.time.Duration

/**
 * Idempotency 키 처리 인터셉터.
 *
 * preHandle: 같은 키로 캐시된 응답이 있으면 컨트롤러 호출 안 하고 그대로 반환.
 * postHandle: 첫 요청이라면 응답을 캐시에 저장.
 *
 * 결제 경로(/payments/**) 에만 적용 — IdempotencyApplication.WebConfig 참조.
 */
@Component
class IdempotencyInterceptor(
    private val store: IdempotencyStore,
) : AsyncHandlerInterceptor {

    private val log = LoggerFactory.getLogger(javaClass)
    private val resTtl = Duration.ofDays(90)  // 결제 — 분쟁 윈도우

    override fun preHandle(req: HttpServletRequest, res: HttpServletResponse, handler: Any): Boolean {
        val key = req.getHeader("Idempotency-Key") ?: return true  // 키 없으면 그냥 통과

        val cached = store.get(key)
        if (cached != null) {
            log.info("idempotency hit key={}", key)
            res.status = cached.status
            res.contentType = "application/json"
            res.writer.write(cached.body)
            res.writer.flush()
            return false  // 컨트롤러 호출 안 함
        }

        // 첫 요청 — 응답을 캡처하기 위해 wrapper 로 교체
        // (request attribute 에 wrapper 저장 후 postHandle 에서 꺼냄)
        val wrapper = CachingResponseWrapper(res)
        req.setAttribute(WRAPPER_ATTR, wrapper)
        return true
    }

    override fun afterCompletion(
        req: HttpServletRequest, res: HttpServletResponse, handler: Any, ex: Exception?,
    ) {
        val key = req.getHeader("Idempotency-Key") ?: return
        val wrapper = req.getAttribute(WRAPPER_ATTR) as? CachingResponseWrapper ?: return
        val body = wrapper.captured()
        // 4xx / 5xx 도 캐시 (재시도 시 같은 결과). 상황에 따라 5xx 만 제외하기도.
        store.put(key, res.status, body, resTtl)
        log.info("idempotency cached key={} status={} ttl_days={}", key, res.status, resTtl.toDays())
    }

    companion object {
        private const val WRAPPER_ATTR = "__idem_wrapper"
    }
}

/**
 * 응답 본문을 캡처하기 위한 wrapper.
 * 운영에서는 Spring 의 ContentCachingResponseWrapper 를 쓰는 게 더 안전.
 */
class CachingResponseWrapper(res: HttpServletResponse) : HttpServletResponseWrapper(res) {
    private val buffer = ByteArrayOutputStream()
    private val writer = PrintWriter(buffer, true, Charsets.UTF_8)
    override fun getWriter(): PrintWriter = writer
    fun captured(): String {
        writer.flush()
        return buffer.toString(Charsets.UTF_8)
    }
}
