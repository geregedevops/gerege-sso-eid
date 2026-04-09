import { prisma } from "./db";

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
  return prisma.orgMember.findUnique({
    where: { organizationId_userId: { organizationId: orgId, userId } },
  });
}

export async function getUserBySubOrThrow(sub: string) {
  const user = await prisma.user.findUnique({ where: { sub } });
  if (!user) throw new Error("User not found");
  return user;
}
