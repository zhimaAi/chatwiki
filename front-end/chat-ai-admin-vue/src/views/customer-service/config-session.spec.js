import { describe, expect, it } from 'vitest'
import { cloneConfig, createSaveQueue, serializeRobotConfig } from './config-session'

const deferred = () => { let resolve; const promise = new Promise(done => { resolve = done }); return { promise, resolve } }

describe('customer configuration saves', () => {
  it('serializes requests and waits for the final submitted save', async () => {
    const queue = createSaveQueue()
    const first = deferred()
    const order = []
    queue.enqueue(async () => { order.push('first'); await first.promise; order.push('first-complete') })
    queue.enqueue(async () => { order.push('second') })
    let flushed = false
    const wait = queue.flush().then(() => { flushed = true })
    await Promise.resolve()
    expect(order).toEqual(['first'])
    expect(flushed).toBe(false)
    first.resolve()
    await wait
    expect(order).toEqual(['first', 'first-complete', 'second'])
  })

  it('blocks testing after a failure and recovers only after a successful retry', async () => {
    const queue = createSaveQueue()
    await expect(queue.enqueue(() => Promise.reject(new Error('save failed')))).rejects.toThrow('save failed')
    await expect(queue.flush()).rejects.toThrow('save failed')
    await queue.enqueue(() => Promise.resolve(true))
    await expect(queue.flush()).resolves.toBeUndefined()
  })

  it('does not mistake a file for an empty JSON object or mutate later snapshots', () => {
    const file = new File(['avatar'], 'avatar.png', { type: 'image/png' })
    const original = { robot_avatar: file, robot_avatar_url: 'old', welcomes: { content: 'welcome', question: [{ content: 'one' }] }, unknown_question_prompt: { content: '', question: ['two'] }, prompt_struct: {}, cache_config: { cache_switch: 1 } }
    const copied = cloneConfig(original)
    const payload = serializeRobotConfig(original)
    original.welcomes.question[0].content = 'changed'
    expect(copied.welcomes.question[0].content).toBe('one')
    expect(payload.robot_avatar).toBe(file)
    expect(payload.robot_avatar_url).toBeUndefined()
    expect(JSON.parse(payload.welcomes).question).toEqual(['one'])
    expect(JSON.parse(payload.unknown_question_prompt).question).toEqual(['two'])
  })
})
