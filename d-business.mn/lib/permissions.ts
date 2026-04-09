import { queryOne } from "./db";

export type OrgRole = "owner" | "admin" | "signer" | "viewer";

const SIGN_ROLES: OrgRole[] = ["owner", "admin", "signer"];
const MANAGE_ROLES: OrgRole[] = ["owner", "admin"];

export function canSign(role: string) {
  return SIGN_ROLES.includes(role as OrgRole);
}

export function canManage(role: string) {
  return MANAGE_ROLES.includes(role as OrgRole);
}

export async function getOrgMembership(userId: string, orgId: string) {
  return queryOne<{ organizationId: string; userId: string; role: string }>(
    `SELECT "organizationId", "userId", role FROM dbiz_org_members WHERE "organizationId"=$1 AND "userId"=$2`,
    [orgId, userId]
  );
}

export async function getUserBySub(sub: string) {
  return queryOne<{ id: string; sub: string; name: string; givenName: string; familyName: string; certSerial: string }>(
    `SELECT id, sub, name, "givenName" as "givenName", "familyName" as "familyName", "certSerial" as "certSerial" FROM dbiz_users WHERE sub=$1`,
    [sub]
  );
}
