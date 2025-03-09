import { cn } from "@/lib/utils";

interface ContainerProps {
  children: React.ReactNode;
  className?: string | undefined;
}

export default function Container({ children, className }: ContainerProps) {
  return (
    <div
      className={cn("mx-auto w-full max-w-screen-xl px-4 md:px-20", className)}
    >
      {children}
    </div>
  );
}
