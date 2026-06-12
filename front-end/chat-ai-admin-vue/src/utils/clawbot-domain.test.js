import { describe, expect, it } from 'vitest'
import {
  buildSetTokenRedirectUrl,
  getCrossDomainTarget,
  isClawbotPath
} from './clawbot-domain'

describe('clawbot-domain utils', () => {
  describe('isClawbotPath', () => {
    it('matches clawbot root path', () => {
      expect(isClawbotPath('/clawbot')).toBe(true)
    })

    it('matches clawbot child path', () => {
      expect(isClawbotPath('/clawbot/chat')).toBe(true)
    })

    it('does not match non-clawbot path', () => {
      expect(isClawbotPath('/robot/list')).toBe(false)
    })
  })

  describe('getCrossDomainTarget', () => {
    it('returns agent domain when entering clawbot', () => {
      expect(
        getCrossDomainTarget({
          toPath: '/clawbot/chat',
          fromPath: '/robot/list',
          agentDomain: 'https://agent.example.com',
          adminDomain: 'https://admin.example.com',
          currentOrigin: 'https://admin.example.com'
        })
      ).toBe('https://agent.example.com')
    })

    it('returns admin domain when leaving clawbot', () => {
      expect(
        getCrossDomainTarget({
          toPath: '/robot/list',
          fromPath: '/clawbot/chat',
          agentDomain: 'https://agent.example.com',
          adminDomain: 'https://admin.example.com',
          currentOrigin: 'https://agent.example.com'
        })
      ).toBe('https://admin.example.com')
    })

    it('does not switch for same module navigation', () => {
      expect(
        getCrossDomainTarget({
          toPath: '/clawbot/settings',
          fromPath: '/clawbot/chat',
          agentDomain: 'https://agent.example.com',
          adminDomain: 'https://admin.example.com',
          currentOrigin: 'https://agent.example.com'
        })
      ).toBe('')
    })

    it('does not switch when target domain is empty', () => {
      expect(
        getCrossDomainTarget({
          toPath: '/clawbot/chat',
          fromPath: '/robot/list',
          agentDomain: '',
          adminDomain: 'https://admin.example.com',
          currentOrigin: 'https://admin.example.com'
        })
      ).toBe('')
    })

    it('does not switch when target domain equals current origin', () => {
      expect(
        getCrossDomainTarget({
          toPath: '/clawbot/chat',
          fromPath: '/robot/list',
          agentDomain: 'https://admin.example.com/',
          adminDomain: 'https://admin.example.com',
          currentOrigin: 'https://admin.example.com'
        })
      ).toBe('')
    })
  })

  describe('buildSetTokenRedirectUrl', () => {
    it('builds a set_token url with encoded redirect path', () => {
      const url = buildSetTokenRedirectUrl({
        domain: 'https://agent.example.com/',
        redirectUrl: '/clawbot/chat?id=1&robot_key=abc',
        token: 'token-1',
        exp: '100',
        ttl: '200',
        userId: '3',
        userName: 'tester'
      })

      expect(url).toBe(
        'https://agent.example.com/#/set_token?token=token-1&exp=100&ttl=200&user_id=3&user_name=tester&redirect_url=%2Fclawbot%2Fchat%3Fid%3D1%26robot_key%3Dabc&refresh_user_info=1'
      )
    })
  })
})
