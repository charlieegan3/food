import { CFP_COOKIE_MAX_AGE } from './constants';
import { sha256, getCookieKeyValue } from './utils';

export async function onRequestPost(context: {
  request: Request;
  env: { CFP_PASSWORD?: string };
}): Promise<Response> {
  const { request, env } = context;
  const body = await request.formData();
  const { password, url } = Object.fromEntries(body);
  const hashedPassword = await sha256(password.toString());
  const hashedCfpPassword = await sha256(env.CFP_PASSWORD);

  if (hashedPassword === hashedCfpPassword) {
    const cookieKeyValue = await getCookieKeyValue(env.CFP_PASSWORD);

    return new Response('', {
      status: 302,
      headers: {
        'Set-Cookie': `${cookieKeyValue}; Max-Age=${CFP_COOKIE_MAX_AGE}; Path=/; HttpOnly; Secure`,
        'Cache-Control': 'no-cache',
        Location: url,
      },
      url: url
    });
  } else {
    // Invalid password. Redirect to login page with error.
    return new Response('', {
      status: 301,
      headers: {
        'Cache-Control': 'no-cache',
        Location: url+'?error=1'
      },
    });
  }
}
